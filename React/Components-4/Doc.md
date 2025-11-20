# Components-3

## Descrição

Este é o terceiro projeto prático do repositório de estudos de React. O diretório `Components-3` foi criado para praticar o uso de lógica JavaScript dentro de componentes React, demonstrando como executar cálculos, operações e declarações de variáveis antes do retorno do componente.

## Estrutura do Projeto

```
Components-3/
├── Component-Friend.jsx  # Componente funcional que usa lógica antes do return
├── App.jsx               # Componente principal que importa e usa ComponentFriend
├── main.jsx              # Ponto de entrada da aplicação React
├── index.html            # Arquivo HTML base
├── vite.config.js        # Configuração do Vite para processar JSX
└── package.json          # Dependências e scripts do projeto
```

## Conceitos Praticados

- **Componentes Funcionais**: Criação de componentes usando funções JavaScript
- **Lógica em Componentes**: Execução de lógica JavaScript antes do retorno JSX
- **Declaração de Variáveis**: Uso de `const`, `let` e operações matemáticas
- **Estruturas de Controle**: Preparação para uso de `if/else`, `switch/case` (fora do return)
- **Exportação/Importação**: Uso de `export default` e `import` para modularizar código
- **JSX**: Sintaxe JavaScript para escrever HTML dentro do código React
- **Arrays e Objetos**: Manipulação de estruturas de dados antes da renderização

## Lógica em Componentes React

### Conceito Fundamental

Em componentes React, toda a lógica JavaScript (cálculos, declarações de variáveis, estruturas condicionais) deve ser executada **antes** do `return`. O `return` deve conter apenas o JSX que será renderizado.

### Exemplo Correto ✅

```jsx
function RandomNumber() {
  // Primeiro, a lógica que deve acontecer antes do retorno
  const n = Math.floor(Math.random() * 10 + 1);

  // Depois, o return usando essa lógica
  return <h1>{n}</h1>;
}
```

**Por que funciona:**
- A variável `n` é declarada e calculada antes do `return`
- O `return` apenas usa o valor já calculado
- A lógica está separada da renderização

### Exemplo Incorreto ❌

```jsx
function RandomNumber() {
  return (
    const n = Math.floor(Math.random() * 10 + 1);
    <h1>{n}</h1>
  )
}
```

**Por que não funciona:**
- Declarações de variáveis (`const`, `let`, `var`) não podem estar dentro do `return`
- O `return` só pode conter expressões JSX, não declarações
- Isso causará um erro de sintaxe

### Regras Importantes

1. **Declarações de variáveis** (`const`, `let`, `var`) devem estar **fora** do `return()`
2. **Estruturas condicionais** (`if/else`, `switch/case`) devem estar **fora** do `return()`
3. **Cálculos e operações** devem ser executados **antes** do `return`
4. O `return()` deve conter apenas **JSX** e **expressões JavaScript** (não declarações)

### Exemplo com Estrutura Condicional

```jsx
function Greeting({ name }) {
  // Lógica antes do return
  let message;
  if (name) {
    message = `Olá, ${name}!`;
  } else {
    message = "Olá, visitante!";
  }

  // Return apenas com JSX
  return <h1>{message}</h1>;
}
```

### Exemplo com Switch/Case

```jsx
function Status({ status }) {
  // Lógica antes do return
  let statusText;
  switch (status) {
    case 'active':
      statusText = 'Ativo';
      break;
    case 'inactive':
      statusText = 'Inativo';
      break;
    default:
      statusText = 'Desconhecido';
  }

  // Return apenas com JSX
  return <p>Status: {statusText}</p>;
}
```

## Instalação e Compilação

### Pré-requisitos

- **Node.js** (versão 18 ou superior)
- **npm** (geralmente vem com o Node.js)

### Comandos para Deixar o Projeto Pronto

**1. Navegue até o diretório do projeto:**
```bash
cd React/Components-3
```

**2. Instale as dependências:**
```bash
npm install
```

Este comando irá:
- Ler o arquivo `package.json`
- Instalar o React e suas dependências
- Instalar o Vite e o plugin do React
- Criar a pasta `node_modules` com todas as dependências
- Gerar o arquivo `package-lock.json` (se ainda não existir)

**3. Inicie o servidor de desenvolvimento:**
```bash
npm run dev
```

O servidor iniciará na porta configurada no `vite.config.js` e abrirá automaticamente no navegador.

**⚠️ IMPORTANTE**: Sempre execute `npm install` primeiro antes de executar `npm run dev`. Se você ver o erro `vite: not found`, significa que as dependências não estão instaladas.

**Nota**: O `package-lock.json` garante que as mesmas versões das dependências sejam instaladas em diferentes ambientes.

### Comandos Disponíveis

1. **Modo de Desenvolvimento** (recomendado para desenvolvimento):
   ```bash
   npm run dev
   ```
   - Inicia o servidor de desenvolvimento do Vite
   - Compila o projeto automaticamente
   - Recarrega a página automaticamente ao fazer alterações (hot reload)
   - O navegador abre automaticamente

2. **Build para Produção** (para criar versão otimizada):
   ```bash
   npm run build
   ```
   - Compila e otimiza o projeto para produção
   - Gera arquivos estáticos na pasta `dist/`
   - Pronto para deploy em servidores web
   - Inclui sourcemaps para debug

3. **Preview da Build de Produção**:
   ```bash
   npm run preview
   ```
   - Pré-visualiza a versão de produção localmente
   - Útil para testar antes de fazer deploy

## Tecnologias Utilizadas

- **React 19.2.0**: Biblioteca JavaScript para construção de interfaces
- **Vite 5.0.0**: Ferramenta de build e servidor de desenvolvimento
- **@vitejs/plugin-react**: Plugin do Vite para processar arquivos JSX

## Componentes

### ComponentFriend
Componente funcional que demonstra o uso de lógica antes do retorno. O componente:
- Declara um array de objetos com dados de amigos
- Seleciona um item específico do array antes do `return`
- Renderiza o título e a imagem do amigo selecionado

### App
Componente principal que serve como container e renderiza o componente `ComponentFriend`.

## Aprendizados

Este projeto demonstra:
- Como executar lógica JavaScript dentro de componentes React
- A importância de separar lógica de renderização
- Como declarar variáveis e executar cálculos antes do `return`
- Como usar estruturas de dados (arrays, objetos) em componentes
- A diferença entre declarações (que vão antes do return) e expressões (que podem ir no JSX)

## Diferenças dos Projetos Anteriores

### Components-1 vs Components-3
- Components-1: Foco em componentes básicos sem lógica complexa
- Components-3: Introduz lógica JavaScript antes do retorno

### Components-2 vs Components-3
- Components-2: Renderização de elementos com propriedades estáticas
- Components-3: Uso de lógica para processar dados antes da renderização

## Erros Comuns a Evitar

⚠️ **NUNCA** coloque declarações (`const`, `let`, `var`) dentro do `return()`

⚠️ **NUNCA** coloque estruturas condicionais (`if/else`, `switch/case`) diretamente no `return()` sem usar operadores ternários ou expressões

✅ **SEMPRE** execute toda a lógica antes do `return`

✅ **SEMPRE** use o `return()` apenas para JSX e expressões JavaScript
