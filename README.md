# Estudo Dev - Repositório Geral

Este repositório centraliza meus projetos práticos, exercícios e anotações de aprendizado em desenvolvimento de software.

## Estrutura do Repositório

### React
Área dedicada ao estudo de **React** com Vite e JavaScript moderno.

- **8 projetos**: Components-1 a Components-5, overview-components, Props-1, Props-2
- Conceitos: componentes funcionais, props, eventos, JSX, manipulação de arrays e listas
- `React/README-REACT.md` — documentação detalhada

### Go
Área dedicada aos estudos práticos de **Go (Golang)**.

- **55+ exercícios** (exercicio1 a exercicio56) que percorrem do básico ao avançado
- Conceitos: sintaxe, variáveis, funções, slices, arrays, structs, ponteiros, métodos, interfaces, pacotes
- `Go/README.md` — lista completa com os conceitos de cada exercício

### GO-API
API RESTful em Go com conexão a PostgreSQL.

- **Task Manager API**: CRUD de tarefas e utilizadores com autenticação JWT, bcrypt para passwords, logs estruturados (slog), tracing (OpenTelemetry) e tratamento de erros consistente
- `GO-API/task_manager/README.md` — documentação completa da API

### GO-POSTGRESQL
Estudo de integração entre **Go e PostgreSQL**.

- 2 exercícios práticos de conexão e execução de queries
- 1 projeto task_manager com migrations SQL
- `GO-POSTGRESQL/README` — detalhes dos exercícios

### Docker
Primeiros passos com **Docker**.

- Clonagem de projeto conteinerizado, modificação de imagem e push para Docker Hub
- `Docker/README-DOCKER.md` — anotações e link do repositório

## Como Executar

Cada projeto tem as suas instruções. Exemplos:

**React:**
```bash
cd React/Components-1
npm install
npm run dev
```

**Go:**
```bash
cd Go
go run hello.go
# ou um exercício específico
cd Go/exercicio1 && go run main.go
```

**Task Manager API:**
```bash
cd GO-API/task_manager
# configurar .env com base no .env.example
go run .
```

## Em Progresso

- **Projetos Full-Stack** — Integrar React (frontend) com Go (backend)

---

*Repositório de estudos pessoais — Aprendizado contínuo em desenvolvimento de software.*
