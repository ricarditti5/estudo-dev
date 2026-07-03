# Go - Exercícios de Aprendizagem

Este diretório contém **55+ exercícios** práticos de Go (Golang), do básico ao avançado.

## Como executar

```bash
cd Go/exercicio1
go run main.go
```

Ou para testar (exercicio56 tem teste):

```bash
cd Go/exercicio56
go test -v
```

## Lista de Exercícios

### 1-10: Sintaxe e Fundamentos

| # | Conceitos |
|---|-----------|
| 1 | Package `main`, import `fmt`, funções personalizadas, `fmt.Println` |
| 2 | Pacotes, exportação de funções entre pacotes |
| 3 | Estrutura de um programa Go |
| 4 | Tipos básicos, declaração de variáveis |
| 5 | Programa mínimo (`package main`, `func main()`, `fmt.Println("Hello")`) |
| 6 | Funções com parâmetros, `Printf`, concatenação de strings |
| 7 | Tipos de dados (`string`, `int`, `bool`, `float64`), zero values |
| 8 | Bloco `var` com múltiplas variáveis |
| 9 | Operações com strings |
| 10 | Bloco `var`, tipos e zero values |

### 11-20: Controlo de Fluxo e Coleções

| # | Conceitos |
|---|-----------|
| 11 | `for` loop, `if/else` |
| 12 | `for` loop com `continue` |
| 13 | `for` loop com `break` |
| 14 | Slices com `make` e `append` |
| 15 | Slices, `append`, reverse iteration com `for` |
| 16 | Iteração sobre slice, índice vs valor |
| 17 | Funções com parâmetros e retorno |
| 18 | Funções com múltiplos valores de retorno |
| 19 | Named returns, `switch` |
| 20 | Arrays (`[5]string`) vs slices (`[]int`), zero values |

### 21-30: Funções e Estruturas de Decisão

| # | Conceitos |
|---|-----------|
| 21 | Funções com parâmetros nomeados e retorno |
| 22 | `if` com initial statement, funções que retornam erro |
| 23 | `if` / `else if` / `else`, comparação de valores |
| 24 | Recursão (função recursiva) |
| 25 | Condicionais `if` / `else if` / `else`, `Printf` com `%d` |
| 26 | `switch` statement |
| 27 | `switch` com fallthrough |
| 28 | `for` loop clássico (inicialização; condição; pós) |
| 29 | `for` como `while` (só condição) |
| 30 | `math/rand`, `switch` com resultado de função |

### 31-40: Structs, Arrays e Ponteiros

| # | Conceitos |
|---|-----------|
| 31 | Arrays (`[n]tipo`) |
| 32 | Slices (`[]tipo`) |
| 33 | Slices com `make` |
| 34 | `range` (índice e valor) |
| 35 | `range` (só valor, ignorar índice com `_`) |
| 36 | `range` (só índice) |
| 37 | Mapas (`map[K]V`) |
| 38 | Mapas com `make` |
| 39 | Mapas — verificar se chave existe (`value, ok := m["key"]`) |
| 40 | Ponteiros (`*int`), dereferenciar para modificar original |

### 41-50: Structs, Métodos e Composição

| # | Conceitos |
|---|-----------|
| 41 | Ponteiros — explicar a diferença entre valor e referência |
| 42 | Structs (`struct`) |
| 43 | Structs com valores atribuídos |
| 44 | Structs aninhadas |
| 45 | Structs — fields nomeados vs posicionais |
| 46 | Métodos em structs (receiver) |
| 47 | Métodos com pointer receiver vs value receiver |
| 48 | Interfaces |
| 49 | Interfaces — implementação implícita |
| 50 | Structs com composição/embedding (`Address`, `Company`) |

### 51-56: Avançado

| # | Conceitos |
|---|-----------|
| 51 | Embedding de structs com métodos |
| 52 | Interfaces com embedding |
| 53 | Type assertions |
| 54 | Type switches |
| 56 | Funções com parâmetros e retorno tipado; testes (`_test.go`) |

## hello.go

Script simples de introdução — primeiro contacto com a sintaxe Go.

## go.mod

Módulo raiz: `my-fisrt-go-project` com Go 1.25.0. Cada exercício tem o seu próprio `go.mod` (module `exercicioN`) — os primeiros usam Go 1.25.0, o exercicio56 usa Go 1.26.2.

---

*Exercícios práticos para aprender Go — do "Hello, World" a interfaces e type assertions.*