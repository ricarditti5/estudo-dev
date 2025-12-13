# Props-1: Entendendo Props e Fluxo de Dados
**Foco:** Passar dados dinâmicos de um componente pai (`App`) para um componente filho (`MyProps`).

## O que são Props?

Cada componente tem algo chamado `props` (abreviação de **propriedades**).

- **Função:** O principal uso de *props* é **passar dados** de um componente pai (superior) para um componente filho (inferior) dentro da hierarquia de componentes. Elas permitem que os componentes se comuniquem e sejam reutilizáveis, de forma semelhante a como passamos argumentos para uma função.

Você já viu isso antes em HTML!
```jsx
<button type="submit" value="Submit"> Submit </button> 
```
Neste exemplo, passamos `type` e `value` para a tag do botão. Da mesma forma, podemos passar informações para nossos próprios componentes!

## Exemplo Prático: `App.jsx` e `MyProps.jsx`

No nosso projeto, veja como estamos passando dados:

1. **Passando a Prop (Componente Pai - `App.jsx`):**
   Estamos passando um atributo chamado `Props` com o valor string `"Hello"`.
   ```jsx
   <MyProps Props="Hello" />
   ```

2. **Recebendo a Prop (Componente Filho - `MyProps.jsx`):**
   Props são um **objeto** JavaScript. O componente filho recebe este objeto contendo todos os atributos passados.
   
   Dentro de `MyProps`, em vez de acessar uma propriedade específica (como `props.Props`), nós convertemos o objeto inteiro para string para visualizá-lo:
   ```jsx
   function MyProps(props) {
       // Converte o objeto { Props: "Hello" } para string
       const stringProps = JSON.stringify(props);
       
       return (
           <div>
               <h1>CHECK OUT MY PROPS OBJECT</h1>
               <h2>{stringProps}</h2>
           </div>
       );
   }
   ```
   **Resultado na tela:** O navegador exibirá `{"Props":"Hello"}`.

## Regras de Props e Sintaxe

- **Imutabilidade:** *Props* são **somente leitura** (*read-only*). O componente `MyProps` recebe os dados, mas não pode alterá-los. Se precisar mudar dados, usaria *state*.

### 1. Passando Strings
Como fizemos no exemplo, se o valor é uma **string**, usamos aspas duplas:
```jsx
// Passa a prop 'Props' com valor "Hello"
<MyProps Props="Hello" />
```
Objeto resultante: `{ Props: "Hello" }`

### 2. Passando Outros Dados (Expressões JavaScript)
Se quiséssemos passar números, booleanos ou variáveis, usaríamos **chaves** `{}`.

Exemplos hipotéticos (se alterássemos o `App.jsx`):
```jsx
// Número
<MyProps age={56} /> 
// Resultado: { age: 56 }

// Arrays
<MyProps items={["A", "B"]} />
// Resultado: { items: ["A", "B"] }
```

## Resumo do Fluxo

1. **App.jsx (Pai):** Define `<MyProps Props="Hello" />`.
2. **React:** Cria o objeto `{ Props: "Hello" }`.
3. **MyProps.jsx (Filho):** Recebe o objeto e (neste caso) exibe-o com `JSON.stringify(props)`.
