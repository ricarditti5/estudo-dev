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

## Instalação e Compilação

### Pré-requisitos

- **Node.js** (versão 16 ou superior)
- **npm** (geralmente vem com o Node.js)

### Criando o package.json

Se você está começando um novo projeto, primeiro precisa criar o arquivo `package.json`:

```bash
npm init -y
```

Este comando cria um `package.json` básico com valores padrão. Em seguida, você precisa instalar as dependências do projeto:

```bash
# Instalar React e React DOM
npm install react react-dom

# Instalar Vite e o plugin do React como dependências de desenvolvimento
npm install -D vite @vitejs/plugin-react
```

### Instalação das Dependências (Projeto Existente)

Se você já possui um `package.json` (como neste projeto), para instalar todas as dependências necessárias, execute:

```bash
npm install
```

Este comando irá:
- Ler o arquivo `package.json`
- Instalar o React e suas dependências
- Instalar o Vite e o plugin do React
- Criar a pasta `node_modules` com todas as dependências
- Gerar o arquivo `package-lock.json` (se ainda não existir)

**Nota**: O `package-lock.json` garante que as mesmas versões das dependências sejam instaladas em diferentes ambientes.

### Compilação e Execução

1. **Modo de Desenvolvimento** (recomendado para desenvolvimento):
   ```bash
   npm run dev
   ```
   - Inicia o servidor de desenvolvimento do Vite
   - Compila o projeto automaticamente
   - Recarrega a página automaticamente ao fazer alterações (hot reload)
   - Acesse a aplicação em `http://localhost:5173`

2. **Build para Produção** (para criar versão otimizada):
   ```bash
   npm run build
   ```
   - Compila e otimiza o projeto para produção
   - Gera arquivos estáticos na pasta `dist/`
   - Pronto para deploy em servidores web

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

