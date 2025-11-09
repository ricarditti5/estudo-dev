# Components-1

## Descrição

Este é o primeiro projeto prático do repositório de estudos de React. O diretório `Components-1` foi criado com o objetivo de praticar os conceitos fundamentais de componentes em React, incluindo criação, importação e renderização de componentes funcionais.

## Estrutura do Projeto

```
Components-1/
├── Component.jsx    # Componente funcional MyComponent
├── App.jsx          # Componente principal que importa e usa MyComponent
├── main.jsx         # Ponto de entrada da aplicação React
├── index.html       # Arquivo HTML base
├── vite.config.js   # Configuração do Vite para processar JSX
└── package.json     # Dependências e scripts do projeto
```

## Conceitos Praticados

- **Componentes Funcionais**: Criação de componentes usando funções JavaScript
- **Exportação/Importação**: Uso de `export default` e `import` para modularizar código
- **JSX**: Sintaxe JavaScript para escrever HTML dentro do código React
- **Props e Estado**: Uso de variáveis dentro de componentes (preparação para props)
- **Renderização**: Uso do `createRoot` do React 18+ para renderizar a aplicação

## Como Executar

1. Instale as dependências (se ainda não instalou):
   ```bash
   npm install
   ```

2. Inicie o servidor de desenvolvimento:
   ```bash
   npm run dev
   ```

3. Acesse a aplicação no navegador (geralmente em `http://localhost:5173`)

## Tecnologias Utilizadas

- **React 19.2.0**: Biblioteca JavaScript para construção de interfaces
- **Vite 5.0.0**: Ferramenta de build e servidor de desenvolvimento
- **@vitejs/plugin-react**: Plugin do Vite para processar arquivos JSX

## Componentes

### MyComponent
Componente funcional simples que exibe uma mensagem de boas-vindas com um nome.

### App
Componente principal que serve como container e renderiza o `MyComponent`.

## Aprendizados

Este projeto demonstra:
- Como estruturar um projeto React básico com Vite
- Como criar e exportar componentes
- Como importar e usar componentes em outros componentes
- Como configurar o Vite para trabalhar com React
- A estrutura básica de uma aplicação React moderna

