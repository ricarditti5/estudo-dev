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

- **Node.js** (versão 18 ou superior)
- **npm** (geralmente vem com o Node.js)

### Comandos para Deixar o Projeto Pronto

**1. Navegue até o diretório do projeto:**
```bash
cd React/Components-1
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

O servidor iniciará na porta **3000** (configurada no `vite.config.js`) e abrirá automaticamente no navegador.

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
   - Acesse a aplicação em `http://localhost:3000`
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

