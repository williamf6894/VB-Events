## Objetivo

Uma plataforma simples de gestão de eventos

### Backend (API)

- CRUD de eventos
  - Nome [Name]
  - Descricao [Description]
  - Data/Hora [Date/Time]
  - Local [Location]
  - Capacidade [Capacity]
  - Duracao [Duration]
- CRUD de participantes
  - Nome [Name]
  - Email 
- Listagem de todos os eventos
  - com filtros por data, status, eventos passados/futuros, etc.
- Verificação de saúde (Health Check)
- Banco de dados (Postgres)

### Frontend

- Lista de todos os eventos
- Detalhes do evento
  - Lista de todos os participantes
  - Cadastro de novo participante
  - Feedback claro sobre:
    - evento lotado
    - entrada inválida
    - etc.

### Tecnologias

- Go (backend)
- Svelte (frontend)
- Postgres (banco de dados)
- README explicativo

### Deploy

- Dockerfile
  - Verificação de saúde (Health Check)
- Docker Compose
  - API
  - Frontend
  - Banco de dados

### Execução

- Instruções no README
- Decisões de arquitetura
