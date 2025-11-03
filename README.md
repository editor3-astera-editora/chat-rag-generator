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

# Observação quanto ao fluxo completo:

Para o fluxo completo, consulte o arquivo `FLOW.md` disponível nesse repositório ou no arquivo `fluxograma.xcalidraw`, que pode ser aberto em: `https://excalidraw.com/`

# Frontend meramente demonstrativo

Um frontend genérico foi montando via LLM somente para demonstração. Para acessá-lo, utilize a pasta web -> `npm install` -> `npm run dev`

# Melhorias futuras

Melhorias futuras que são necessárias estão no arquivo `MELHORIAS_FUTURAS.md` disponível nesse repositório.

