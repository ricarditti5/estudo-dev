## Anotações
Event listener em componentes

## O que aprendi neste diretório: Event Listeners

Neste módulo, o foco foi aprender a adicionar interatividade aos componentes React através de **Event Listeners** (ouvintes de eventos).

### Conceitos Principais:

1.  **Atributos de Evento:**
    - O React utiliza atributos camelCase para eventos, como `onClick` em vez de `onclick` do HTML padrão.
    - Exemplo: `<button onClick={...}>`

2.  **Passando Funções:**
    - Passamos uma *função* para o manipulador de eventos, não uma string.
    - A função define o que acontece quando o evento ocorre (ex: clicar no botão).

### Exemplo do Código (`Button.jsx`):

```jsx
function Button(){
    // Função que será executada ao clicar
    function MyBotton(){
        Alert('This is a button');
    }
    
    // Passando a função MyBotton para o evento onClick
    return <button onClick={MyBotton}>Click here</button>
}
```

**Resumo:** Aprendi a fazer um botão "fazer algo" quando clicado, definindo uma função dentro do componente e conectando-a ao evento `onClick` do elemento JSX.
