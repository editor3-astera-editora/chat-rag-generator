`/cmd/main.go`

- Função: ponto de entrada da aplicação.
- Responsabilidades:
    - Inicializa a conexão com  o banco de dados (config.GetDB())
    - Cria instâncias do repositório, serviço e handler.
    - Configura o roteamento HTTP (/chat)
    - Aplica o middleware de logging
    - Inicia o servidor local em `http://localhost:8080`

`/internal/handler/chat_handler.go`

- Função: Controlador HTTP responsável por receber e responder às requisições REST do chat.
- Responsabilidades:
    - Recebe mensagens do usuário via `POST /chat`.
    - Decodifica o corpo JSON (`user_id`, `message`)
    - Encaminha para o `ChatService`
    - Retorna a resposta do modelo como JSON

Exemplo de entrada:

```
{
    "user_id": "user123".
    "message": "Explique a leia de Ohm."
}
```

Saída esperada:

```
{
    "response": "A lei de Ohm afirma que V = R x I."
}
```

`/internal/llm/embedding.go`

- Função: Gera embeddings vetoriais das mensagens.
- Responsabilidades: 
    - Utiliza o endpoint de embeddings da API OpenAI (`text-embedding-3-large`).
    - Converte o vetor retornado (float64) para []float32 para compatibilidade com o PostGreSQL (pgVector)
    - Serve de base para a busca vetorial no banco

`/internal/llm/openai_client.go`

- Função: faz chamadas ao modelo de linguagem da OpenAI
- Responsabilidades:
    - utiliza o modelo gpt-4o para gerar respostas com base no contexto.
    - Define o prompt do sistema como um "professor que explica conceitos de livros didáticos."
    - Recebe o conteúdo completo (mensagem + contexto + histórico) e retorna a resposta textual

`/internal/memory/local_memory.go`

- Função: implementa a memória local por usuário.
- Responsabildiades:
    - Armazena mensagens em um map[UserID][]Message
    - Mantém apenas as últimas 10 mensagens (MaxMessages)
    - Permite recuperar ou limpar o histórico do usuário

Limite configurado: MaxMessages = 10

`/internal/middleware/logging.go`

- Função: middleware para logar requisições HTTP
- Responsabilidades:
    - Registra o método, rota e tempo de execução de cada requisição.
    - Auxilia no monitoramento e debugging do servidor

Saída típica: 
    POST /chat
    POST /chat em 132ms

`/internal/model/mesage.go`

- Função: define a estrutura das mensagens trocadas entre usuário e assistente.
- Campos:
    - Role: "user" ou "assistant"
    - Content: texto da mensagem 

`/internal/repository/chat_repository.go`

- Função: interage com o banco de dados vetorial (pgVector)
- Responsabilidades:
    - Executa busca semântica no banco usando similaridade vetorial (<->)
    - Retorna os documentos mais próximos do embedding da pergunta.
    - Limita o número de resultados relevantes (LIMIT configurável)
- Consulta SQL:

```
SELECT document
FROM langchain_pg_embedding
ORDER BY embedding <-> $1
LIMIT $2
```

`/internal/service/chat_service.go`
- Função: Camada de orquestração principal do chat
- Fluxo interno:
    1. Armazena a nova mensagem do usuário na memória
    2. Gera embedding do conteúdo via  `llm.GenerateEmbedding`
    3. Busca documentos similares no banco via `ChatRepository`
    4. Constrói um contexto concatenando resultados relevantes e histórico.
    5. Chama `llm.GetLLMResponse()` com o contexto completo.
    6. Armazena a respost ado modelo na memória.
    7. Retorna a resposta final


    A[Usuário envia mensagem] --> B[Handler (/chat)]
    B --> C[ChatService]
    C --> D[GenerateEmbedding]
    D --> E[pgvector via ChatRepository]
    E --> F[Busca contexto relevante]
    F --> G[Compõe histórico + contexto]
    G --> H[GetLLMResponse (GPT-4o)]
    H --> I[Memória atualizada (últimas 10 msgs)]
    I --> J[Resposta enviada ao usuário]
