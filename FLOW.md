```mermaid
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
