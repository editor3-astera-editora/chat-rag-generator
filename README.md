# Como iniciar o projeto

Esse projeto implementa o MVP de um backend de um professor virtual que utiliza um pipeline RAG (retrieval-augmented generation) conectado ao banco vetorial pgVector, alimentado pelo pipeline `https://github.com/editor3-astera-editora/rag-pgvector`.
Ele permite que o chatbot recupere informações dos livros didáticos, gere respostas contextuais e explique fórmulas matemáticas de forma pegagógica.

## Estrutura do projeto 

```
cmd/main.go                              → Ponto de entrada do servidor HTTP
internal/
 ├── handler/chat_handler.go             → Controla requisições HTTP (camada de interface)
 ├── service/chat_service.go             → Orquestra o fluxo de inferência e memória
 ├── repository/chat_repository.go       → Acesso a embeddings (pgVector)
 ├── repository/formula_repository.go    → Acesso ao mapa de fórmulas
 ├── llm/embedding.go                    → Geração de embeddings com OpenAI
 ├── llm/openai_client.go                → Chamada de completions (GPT-4o)
 ├── memory/local_memory.go              → Armazena o histórico de conversa (por usuário)
 ├── middleware/                         → Middlewares de logging, CORS e checagens semânticas
 └── model/                              → Estruturas de dados (Message, Formula)
```

Visão geral do fluxo do projeto:

```
HTTP → handler → service → (llm + repository + memory) → handler → HTTP
                   └─ middleware de apoio (CORS, logging, checagens)
```

## Inicialização:

1. Variáveis de ambiente

Crie um arquivo `.env`com:

```
OPENAI_API_KEY="s..."
PGVECTOR_DB_URI=postgresql+psycopg://postgres:SUASENHA@localhost:5432/embeddings_db
PSYCOPG_DB_URI=postgresql://postgres:SUASENHA@localhost:5432/embeddings_db
````

O banco `embeddings_db` deve ter sido criado com o schema do projeto `https://github.com/editor3-astera-editora/rag-pgvector`

2. Pré-requisitos

É necessário ter instalado na máquina Golang - 1.23.3 (mínimo). O download pode ser realizado em: `https://go.dev/doc/install`

3. Execute `go mod tidy` para instalar as bibliotecas necessárias.

4. Execute `go run ./cmd`:

Saída esperada:

```
 Conexão PostgreSQL (pgxpool) bem-sucedida
2025/11/03 17:25:32  Banco conectado: embeddings_db | Usuário: postgres | Versão: 16.10
2025/11/03 17:25:32  search_path: "$user", public
2025/11/03 17:25:32 Servidor inicado em http://localhost:8080
```

## API 

Rota principal:

```
POST /chat
```
Entrada esperada (JSON):

```
{
  "user_id": "user123",
  "message": "Como calcular o montante em juros compostos?"
}
```

Teste com:

```
Invoke-RestMethod -Uri "http://localhost:8080/chat" `
>>   -Method Post `
>>   -Body '{"user_id":"user123","message":"Como calcular juros compostos?"}' `
>>   -ContentType "application/json"
```

## Processo interno:

1. O `ChatHandler` decodifica o JSON recebido;
2. O `ChatService` gera o embedding da pergunta (`text-embedding-3-large`);
3. O `ChatRepository` busca no banco pgvector os 3 chunks mais similares;
4. Se a similaridade média for menor que 0.75, o sistema responde com uma mensagem de redirecionamento pedagógico;
5. Caso contrário:
       - monta o **contexto textual** a partir dos trechos recuperados;
       - busca **fórmulas associadas** ao capítulo (`FormulaService` + `formulas_map`);
       - adiciona memória conversacional (`local_memory.go`);
       - envia tudo ao modelo `gpt-4o` para gerar a resposta final.

## Saída (JSON):
```
{
  "response": "Para calcular o montante em juros compostos, usamos a fórmula M = C × (1 + i)^t. Esse conceito está explicado no capítulo 3 da Unidade 2.",
  "sources": [
    {
      "Document": "O montante é o valor final de uma aplicação financeira...",
      "BookName": "Matemática Financeira",
      "Unit": 2,
      "Chapter": 3,
      "Similarity": 0.89
    }
  ]
}
```


```
flowchart TD

A[POST /chat recebido]:::handler --> B{AuthMiddleware<br>(usuário autorizado?)}:::middleware
B -->|Não| E[403 Forbidden<br>→ encerra]:::middleware
B -->|Sim| F{CORSMiddleware<br>(origem permitida?)}:::middleware
F -->|Não| G[403 CORS Bloqueado]:::middleware
F -->|Sim| H{Método == POST?}:::handler
H -->|Não| I[405 Method Not Allowed]:::handler
H -->|Sim| J{JSON válido?}:::handler
J -->|Não| K[400 Bad Request<br>Erro ao decodificar JSON]:::handler
J -->|Sim| L[ChatHandler chama<br>ChatService.ProcessMessage()]:::handler

%% SERVICE LAYER
L --> M{userBooks[userID] existe?}:::service
M -->|Não| N[404 Livro não mapeado]:::service
M -->|Sim| O{Pergunta contém<br>'capítulo X'?}:::service
O -->|Sim| P[Define chapterFilter = X]:::service
O -->|Não| Q[chapterFilter = nil]:::service
P --> R[Gerar embedding<br>(text-embedding-3-large)]:::llm
Q --> R
R -->|Erro| S[500 Erro ao gerar embedding]:::llm
R -->|OK| T[Repository.SearchSimilarEmbeddings()]:::repository

%% REPOSITORY
T -->|Nenhum resultado| U[Sem contexto relevante<br>→ Resposta pedagógica padrão]:::service
T -->|Com resultados| V{Maior Similaridade >= 0.75?}:::service
V -->|Não| U
V -->|Sim| W[Construir contexto textual<br>com docs recuperados]:::service

%% FORMULA LOOKUP
W --> X{Pergunta sugere fórmula?<br>('como calcular', 'fórmula', ...)}:::middleware
X -->|Não| Y[Sem fórmulas adicionadas]:::service
X -->|Sim| Z[Repository.SearchFormula()<br>no formulas_map]:::repository
Z -->|Encontrou| Z1[Adiciona blocos de fórmulas ao contexto]:::service
Z -->|Nenhuma| Y
Z1 --> AA[Montar prompt final<br>com histórico + contexto + pergunta]:::service
Y --> AA

%% LLM GENERATION
AA --> AB[LLM.GetLLMResponse()<br>(GPT-4o)]:::llm
AB -->|Erro| AC[500 Erro ao chamar LLM]:::llm
AB -->|OK| AD[AddMessage(userID, resposta)]:::memory
AD --> AE[Montar JSON final<br>{response, sources}]:::handler
AE --> AF[Retornar resposta ao cliente<br>HTTP 200 OK]:::handler

%% END
AF --> AG[Fim do fluxo]:::end

%% STYLES
classDef handler fill:#2f4f4f,stroke:#111,color:#fff
classDef service fill:#1f77b4,stroke:#111,color:#fff
classDef repository fill:#9467bd,stroke:#111,color:#fff
classDef middleware fill:#ff7f0e,stroke:#111,color:#fff
classDef llm fill:#2ca02c,stroke:#111,color:#fff
classDef memory fill:#8c564b,stroke:#111,color:#fff
classDef end fill:#333,stroke:#000,color:#fff
```
