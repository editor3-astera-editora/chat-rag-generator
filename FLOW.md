```mermaid
flowchart TD

A[POST /chat recebido]:::handler --> B{AuthMiddleware: usuario autorizado?}:::middleware
B -->|Nao| E[403 Forbidden - encerra]:::middleware
B -->|Sim| F{CORSMiddleware: origem permitida?}:::middleware
F -->|Nao| G[403 CORS bloqueado]:::middleware
F -->|Sim| H{Metodo POST?}:::handler
H -->|Nao| I[405 Method Not Allowed]:::handler
H -->|Sim| J{JSON valido?}:::handler
J -->|Nao| K[400 Bad Request - Erro JSON]:::handler
J -->|Sim| L[Handler chama ChatService ProcessMessage]:::handler

%% SERVICE LAYER

L --> M{Usuario possui livro associado?}:::service
M -->|Nao| N[404 Livro nao mapeado]:::service
M -->|Sim| O{Pergunta contem capitulo X?}:::service
O -->|Sim| P[Define chapterFilter = X]:::service
O -->|Nao| Q[chapterFilter = nil]:::service
P --> R[Gerar embedding text-embedding-3-large]:::llm
Q --> R
R -->|Erro| S[500 Erro ao gerar embedding]:::llm
R -->|OK| T[Repository SearchSimilarEmbeddings]:::repository

%% REPOSITORY

T -->|Nenhum resultado| U[Sem contexto relevante -> resposta pedagogica generica]:::service
T -->|Com resultados| V{Maior similaridade >= 0.75?}:::service
V -->|Nao| U
V -->|Sim| W[Construir contexto textual com docs recuperados]:::service

%% FORMULA LOOKUP

W --> X{Pergunta sugere formula?}:::middleware
X -->|Nao| Y[Sem formulas adicionadas]:::service
X -->|Sim| Z[Repository SearchFormula no formulas_map]:::repository
Z -->|Encontrou| Z1[Adiciona blocos de formulas ao contexto]:::service
Z -->|Nenhuma| Y
Z1 --> AA[Montar prompt final com historico e pergunta]:::service
Y --> AA

%% LLM GENERATION

AA --> AB[LLM GetLLMResponse GPT-4o]:::llm
AB -->|Erro| AC[500 Erro ao chamar LLM]:::llm
AB -->|OK| AD[AddMessage na memoria local]:::memory
AD --> AE[Montar JSON final com response e sources]:::handler
AE --> AF[Retornar resposta HTTP 200 OK]:::handler
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
