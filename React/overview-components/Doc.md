# Visão Geral dos Componentes React

Este documento serve como um guia de estudo centralizado, explicando os conceitos abordados em cada diretório de componentes.

## Components-1: Componentes Básicos e Variáveis
**Foco:** Criação de um componente funcional simples e uso de variáveis locais.

Neste exemplo, temos um componente `MyComponent` que define uma variável `nome` e a exibe dentro de um JSX.
- **Conceito Chave:** Interpolação de variáveis no JSX usando chaves `{}`.
- **Código:**
  ```jsx
  export default function MyComponent() {
      const nome = 'Ricardo';
      return (
          <>
              <h1>Olá {nome}</h1>
          </>
      );
  }
  ```

## Components-2: Objetos e Propriedades
**Foco:** Uso de objetos para armazenar dados e acessá-los no JSX.

Aqui, os dados da coruja (título e imagem) são armazenados em um objeto `owl` fora do componente. O componente então acessa essas propriedades.
- **Conceito Chave:** Acesso a propriedades de objetos (`owl.title`, `owl.src`) para renderizar conteúdo dinâmico e atributos de tags HTML (como `src` e `alt` em `<img>`).
- **Código:**
  ```jsx
  const owl = {
      title: 'Excelent owl',
      src: 'https://content.codecademy.com/courses/React/react_photo-owl.jpg',
  };
  // ... dentro do componente:
  // <h1>{owl.title}</h1>
  // <img src={owl.src} alt={owl.title} />
  ```

## Components-3 & Components-4: Arrays e Listas
**Foco:** Manipulação de Arrays de objetos.

Nestes diretórios, temos uma lista de amigos (`friends`), que é um array de objetos. O componente seleciona um item específico desse array para renderizar.
- **Conceito Chave:** Acesso a elementos de um array pelo índice (`friends[1]`) e renderização dos dados desse item.
- **Código:**
  ```jsx
  const friends = [
    { title: "Yummmmmmm", src: "..." },
    { title: "Hey Guys! Wait Up!", src: "..." },
    // ...
  ];
  
  function ComponentFriend(){
      const friend = friends[1]; // Seleciona o segundo amigo
      return(
          <>
              <h1>{friend.title}</h1>
              <img src={friend.src} />
          </>
      )
  }
  ```

## Components-5: Eventos (Event Handlers)
**Foco:** Interatividade e manipulação de eventos.

Este exemplo introduz um botão que dispara uma ação quando clicado.
- **Conceito Chave:** Adicionar ouvintes de eventos (`onClick`) a elementos JSX e passar uma função de callback (`MyBotton`) para ser executada quando o evento ocorrer.
- **Código:**
  ```jsx
  function Button(){
      function MyBotton(){
          Alert('This is a button');
      }
      return <button onClick={MyBotton}>Click here</button>
  }
  ```
