# Components-2

## Descrição

Este é o segundo projeto prático do repositório de estudos de React. O diretório `Components-2` foi criado para praticar componentes que renderizam elementos com propriedades, como imagens e textos, demonstrando o uso de objetos JavaScript dentro de componentes React.

## Estrutura do Projeto

```
Components-2/
├── Coruja.jsx      # Componente funcional Coruja que exibe uma imagem
├── App.jsx         # Componente principal que importa e usa Coruja
├── main.jsx        # Ponto de entrada da aplicação React
├── index.html      # Arquivo HTML base
├── vite.config.js  # Configuração do Vite para processar JSX
└── package.json    # Dependências e scripts do projeto
```

## Conceitos Praticados

- **Componentes Funcionais**: Criação de componentes usando funções JavaScript
- **Exportação/Importação**: Uso de `export default` e `import` para modularizar código
- **JSX**: Sintaxe JavaScript para escrever HTML dentro do código React
- **Objetos JavaScript**: Uso de objetos para armazenar dados (como URLs de imagens)
- **Atributos JSX**: Uso de atributos como `src` e `alt` em elementos HTML
- **Renderização**: Uso do `createRoot` do React 18 para renderizar a aplicação

## Instalação e Compilação

### Pré-requisitos

- **Node.js** (versão 18 ou superior)
- **npm** (geralmente vem com o Node.js)

### Comandos para Deixar o Projeto Pronto

**1. Navegue até o diretório do projeto:**
```bash
cd React/Components-2
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

O servidor iniciará na porta **3001** (configurada no `vite.config.js`) e abrirá automaticamente no navegador.

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
   - Acesse a aplicação em `http://localhost:3001`
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

### Coruja
Componente funcional que exibe uma imagem de coruja usando dados de um objeto JavaScript. O componente demonstra:
- Uso de objetos para armazenar dados
- Renderização de imagens com atributos `src` e `alt`
- Interpolação de valores JavaScript no JSX

### App
Componente principal que serve como container e renderiza o componente `Coruja`.

## Aprendizados

Este projeto demonstra:
- Como estruturar um projeto React básico com Vite
- Como criar e exportar componentes
- Como importar e usar componentes em outros componentes
- Como usar objetos JavaScript dentro de componentes
- Como renderizar elementos HTML com atributos dinâmicos
- Como configurar o Vite para trabalhar com React
- A estrutura básica de uma aplicação React moderna

## Diferenças do Components-1

- Este projeto foca em renderização de elementos com propriedades (imagens)
- Demonstra o uso de objetos para organizar dados
- Mostra como usar atributos HTML em JSX de forma dinâmica

