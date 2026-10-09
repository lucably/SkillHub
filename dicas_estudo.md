# 📘 Anotações de Go — Backend Development

> Material de consulta sobre módulos, ponteiros, interfaces, contextos, goroutines e channels em Go.

---

## 📦 1. Gerenciamento de dependências

### `go.mod` — equivalente aproximado ao `package.json`

No ecossistema JavaScript, usamos o `package.json` para definir informações do projeto e suas dependências.

Em Go, utilizamos o arquivo `go.mod` para declarar o módulo e gerenciar suas dependências.

| Go | JavaScript / Node.js |
|---|---|
| `go.mod` | `package.json` |
| `go get github.com/gin-gonic/gin` | `npm install express` |
| `go.sum` | Sem equivalente exato; lembra o papel de integridade de arquivos de lock |

**Atenção:** as equivalências são conceituais. `go.mod` e `package.json` possuem diferenças de funcionamento, e `go.sum` não é exatamente equivalente ao `package-lock.json`.

### O que é `go.sum`?

Pense nele como um arquivo que registra checksums criptográficos utilizados para verificar a integridade do conteúdo das dependências.

- **`go.mod`** → declara o módulo e os requisitos de dependências.
- **`go.sum`** → registra checksums para verificar a integridade de versões de módulos baixadas.

### O que é `go mod tidy`?

O comando `go mod tidy` analisa os imports do projeto e sincroniza as dependências declaradas nos arquivos do módulo.

Por exemplo, começamos a utilizar o Gin:

```go
import "github.com/gin-gonic/gin"
```

Depois, executamos:

```bash
go mod tidy
```

O Go ajusta as dependências necessárias, adicionando requisitos que estejam faltando e removendo aqueles que não são mais necessários.

**Resumo:** pense no `go mod tidy` como uma ferramenta de limpeza e sincronização das dependências do projeto.

---

## 🧱 2. Structs e valores padrão

Quando inicializamos uma struct vazia, o Go atribui automaticamente o valor zero correspondente ao tipo de cada campo.

Quando utilizamos ponteiros, o valor zero do ponteiro é `nil`.

### Exemplo

```go
type EstruturaValor struct {
    Ativo bool
    Nome  string
    Idade int
}

type EstruturaPonteiro struct {
    Ativo *bool
    Nome  *string
    Idade *int
}

func main() {
    v := EstruturaValor{}

    fmt.Printf(
        "Valor -> Ativo: %t, Nome: %q, Idade: %d\n",
        v.Ativo, v.Nome, v.Idade,
    )

    // Saída: Valor -> Ativo: false, Nome: "", Idade: 0

    p := EstruturaPonteiro{}

    fmt.Printf(
        "Ponteiro -> Ativo: %v, Nome: %v, Idade: %v\n",
        p.Ativo, p.Nome, p.Idade,
    )

    // Saída: Ponteiro -> Ativo: <nil>, Nome: <nil>, Idade: <nil>
}
```

### Valores zero mais comuns

| Tipo | Valor zero |
|---|---|
| `bool` | `false` |
| `string` | `""` |
| `int` | `0` |
| `float64` | `0` |
| Ponteiros | `nil` |
| Slices | `nil` |
| Maps | `nil` |
| Interfaces | `nil`, quando não possuem valor dinâmico nem tipo dinâmico |

### Quando utilizar ponteiros?

Uma regra melhor do que simplesmente dizer "precisa alterar dados? Use ponteiro" é avaliar o que a função precisa fazer.

- **Valor:** recebe uma cópia do valor. Alterações feitas nessa cópia não modificam a variável original.
- **Ponteiro:** recebe um endereço que permite acessar o valor original e modificá-lo.
- **Ponteiro também pode ser útil:** para evitar cópias de estruturas maiores ou representar a ausência de um valor.

Exemplo:

```go
type User struct {
    Name string
}

func changeName(user User) {
    user.Name = "Maria"
}

func changeNameWithPointer(user *User) {
    user.Name = "Maria"
}

func main() {
    user := User{Name: "Lucas"}

    changeName(user)
    fmt.Println(user.Name) // Lucas

    changeNameWithPointer(&user)
    fmt.Println(user.Name) // Maria
}
```

**Importante:** a cópia não significa que todos os dados serão perdidos ao final da execução. Significa que as alterações feitas na cópia não modificam diretamente a variável original. Em tipos como slices e maps, existem particularidades porque eles referenciam estruturas de dados compartilhadas.

---

## 🔒 3. Escopo de variáveis no `if`

Em Go, podemos declarar variáveis diretamente na inicialização de um `if`.

Essas variáveis ficam disponíveis no próprio `if`, nos seus blocos `else if` e `else`, mas não fora dessa estrutura.

### Exemplo

```go
func buscarIdade() (int, error) {
    return 25, nil
}

func main() {
    if idade, err := buscarIdade(); err != nil {
        fmt.Println("Erro ao buscar idade:", err)
    } else {
        fmt.Printf("A idade é %d\n", idade)
    }

    // Erro de compilação:
    // fmt.Println(idade)
    // undefined: idade
}
```

Observe:

```go
if idade, err := buscarIdade(); err != nil {
    // ...
}
```

A expressão antes do ponto e vírgula declara as variáveis `idade` e `err`. Elas ficam restritas ao escopo dessa estrutura condicional.

Essa abordagem é bastante comum em Go para tratar erros sem criar variáveis desnecessárias fora do bloco.

---

## 🔌 4. Interfaces em Go

Uma interface define um conjunto de métodos que um tipo precisa possuir para satisfazê-la.

Em Go, a implementação de interfaces é implícita: não precisamos declarar `implements`, como ocorre em algumas outras linguagens.

### Exemplo: pagamentos

```go
type PaymentProcessor interface {
    Pay(amount float64) error
}

type PixPayment struct{}

type CreditCardPayment struct{}

func (p PixPayment) Pay(amount float64) error {
    fmt.Println("Processing PIX payment:", amount)
    return nil
}

func (c CreditCardPayment) Pay(amount float64) error {
    fmt.Println("Processing Credit Card payment:", amount)
    return nil
}

func processPayment(
    payment PaymentProcessor,
    amount float64,
) error {
    return payment.Pay(amount)
}

func main() {
    _ = processPayment(PixPayment{}, 100)
    _ = processPayment(CreditCardPayment{}, 200)

    // Também podemos chamar diretamente:
    _ = PixPayment{}.Pay(100)
}
```

### Como isso funciona?

A interface `PaymentProcessor` exige o método:

```go
Pay(amount float64) error
```

Qualquer tipo que possua esse método com a assinatura exata satisfaz a interface.

No exemplo, tanto `PixPayment` quanto `CreditCardPayment` implementam `PaymentProcessor`.

A função `processPayment` recebe uma interface, não uma implementação específica.

```go
func processPayment(
    payment PaymentProcessor,
    amount float64,
) error {
    return payment.Pay(amount)
}
```

Isso permite passar diferentes meios de pagamento sem precisar criar uma função diferente para cada um.

### Tipo concreto versus interface

Podemos chamar diretamente:

```go
PixPayment{}.Pay(100)
```

Nesse caso, estamos utilizando o tipo concreto `PixPayment`.

Por outro lado:

```go
processPayment(PixPayment{}, 100)
```

Estamos passando `PixPayment` para uma função que aceita qualquer valor que satisfaça `PaymentProcessor`.

**Vantagem:** podemos adicionar novos meios de pagamento sem precisar modificar a lógica central de processamento.

### E se adicionarmos outro método à interface?

Imagine que a interface passe a exigir `Show()`:

```go
type PaymentProcessor interface {
    Pay(amount float64) error
    Show()
}
```

Agora, `PixPayment` precisa implementar os dois métodos.

```go
func (p PixPayment) Pay(amount float64) error {
    fmt.Println("PIX:", amount)
    return nil
}

func (p PixPayment) Show() {
    fmt.Println("Showing PIX payment")
}
```

Se `Show()` estiver ausente, `PixPayment` não satisfará a nova interface.

Uma chamada como:

```go
processPayment(PixPayment{}, 100)
```

não compilará se `processPayment` receber a interface atualizada, pois `PixPayment` não implementará todos os métodos exigidos.

**Dica:** interfaces pequenas costumam ser mais fáceis de implementar, testar e reutilizar.

---

## ⏱️ 5. O que é `context.Context`?

Simplificando, `context.Context` carrega informações sobre o ciclo de vida de uma operação.

Suas principais funcionalidades incluem:

- Cancelamento de operações.
- Timeout e deadline.
- Valores associados ao contexto de uma operação.

No desenvolvimento backend, cancelamento e timeout são especialmente importantes.

### Criando um contexto raiz

```go
func main() {
    ctx := context.Background()

    fmt.Println(ctx)
}
```

`context.Background()` cria um contexto vazio, sem cancelamento ou deadline.

Pense nele como o ponto de partida de uma operação.

Entretanto, em uma API real, o contexto geralmente vem de uma requisição HTTP ou de uma operação anterior e é propagado pelas camadas da aplicação.

### Analogia com o frontend

Imagine que o frontend inicia uma requisição HTTP e o usuário sai da tela antes de ela terminar.

No frontend, podemos cancelar uma requisição usando `AbortController`.

No backend Go, `context.Context` permite propagar sinais de cancelamento e prazos para operações relacionadas, desde que as funções envolvidas respeitem o contexto.

---

## 🛑 6. Contexto, `select` e timeout

O `select` permite esperar por operações em canais.

Quando vários casos estão prontos ao mesmo tempo, o Go escolhe um deles. Portanto, não é correto afirmar que ele sempre seleciona o evento que ocorreu primeiro cronologicamente.

### Exemplo

```go
func doSomething(ctx context.Context) error {
    select {
    case <-ctx.Done():
        fmt.Println("Operação foi cancelada.")
        return ctx.Err()

    case <-time.After(5 * time.Second):
        fmt.Println("Operação terminou.")
        return nil
    }
}

func main() {
    ctx, cancel := context.WithTimeout(
        context.Background(),
        2*time.Second,
    )
    defer cancel()

    err := doSomething(ctx)
    fmt.Println("Erro:", err)
}
```

Nesse exemplo, o contexto tem timeout de dois segundos, enquanto a operação simulada leva cinco segundos.

O resultado esperado é:

```text
Operação foi cancelada.
Erro: context deadline exceeded
```

### Por que isso acontece?

O `context.WithTimeout` cria um contexto que será cancelado quando o prazo expirar.

```go
ctx, cancel := context.WithTimeout(
    context.Background(),
    2*time.Second,
)
```

A operação espera por dois possíveis eventos:

1. O contexto ser cancelado.
2. O timer de cinco segundos terminar.

Como o timeout do contexto ocorre antes, o caso abaixo é selecionado:

```go
case <-ctx.Done():
```

Se o contexto continuar ativo e o timer de cinco segundos terminar primeiro, a função retornará `nil`.

### O que significa `<-ctx.Done()`?

O método `Done()` retorna um canal que é fechado quando o contexto é cancelado.

```text
ctx.Done()
    ↓
Canal sinalizado quando o contexto é cancelado
```

O operador `<-` permite receber de um canal. Quando usado dentro de um `select`, o programa espera até que algum dos casos esteja pronto.

**Atenção:** `ctx.Err()` retorna o motivo do cancelamento, como `context.Canceled` ou `context.DeadlineExceeded`. Se a operação terminar normalmente, retorne `nil`, e não `ctx.Err()`.

---

## 🧵 7. Goroutines e channels

### O que é uma goroutine?

Uma goroutine permite executar uma função concorrentemente.

Por exemplo:

```go
go doSomething()
```

Isso inicia a função em uma nova goroutine, permitindo que a goroutine atual continue sem esperar que ela termine.

```text
main
 │
 ├── go doSomething()
 │       ↓
 │    executando concorrentemente
 │
 └── continua a execução
```

O início da goroutine não garante que ela será executada imediatamente, nem que continuará executando depois que a função `main` terminar.

### Concorrência versus paralelismo

- **Concorrência:** organizar e executar várias tarefas de forma que seu progresso possa se intercalar.
- **Paralelismo:** executar tarefas literalmente ao mesmo tempo em diferentes recursos de processamento.

Go oferece suporte a ambos. O runtime administra as goroutines e sua execução.

### O que é um channel?

Um channel é um mecanismo de comunicação e sincronização entre goroutines.

```go
ch := make(chan string)
```

Nesse exemplo, criamos um canal que transporta valores do tipo `string`.

### Exemplo completo

```go
func sayHello(ch chan string) {
    time.Sleep(2 * time.Second)

    ch <- "Hello"
}

func main() {
    ch := make(chan string)

    go sayHello(ch)

    message := <-ch

    fmt.Println(message)
}
```

Saída:

```text
Hello
```

### Entendendo a execução

```text
main
 │
 ├── cria o channel
 │
 ├── inicia a goroutine
 │       │
 │       └── espera 2 segundos
 │
 ├── aguarda receber de ch
 │
 │       2 segundos depois
 │       ↓
 │    ch <- "Hello"
 │
 └── recebe "Hello" e imprime
```

Quando executamos:

```go
message := <-ch
```

a goroutine atual fica bloqueada até conseguir receber um valor do canal.

```text
<-ch
 │
 ├── Existe um valor disponível?
 │       └── Sim: recebe e continua
 │
 └── Não existe um valor disponível?
         └── Espera até conseguir receber
```

Como o channel não tem buffer, o envio e o recebimento precisam se sincronizar.

### O que é `make`?

`make` é uma função embutida de Go usada para inicializar slices, maps e channels.

Exemplos:

```go
slice := make([]int, 3)
m := make(map[string]int)
ch := make(chan string)
```

Declarar essas estruturas sem inicializá-las pode resultar em um valor `nil`.

Por exemplo:

```go
var ch chan string

fmt.Println(ch == nil) // true
```

Um channel `nil` não permite comunicação: envios e recebimentos ficam bloqueados indefinidamente, e um `select` que só contenha operações em canais `nil` também pode ficar bloqueado.

---

## 📬 8. Channels com buffer

Um channel com buffer permite armazenar uma quantidade limitada de valores antes de precisar sincronizar o envio com um recebimento.

```go
ch := make(chan string, 1)
```

O número `1` indica que o canal consegue armazenar um valor sem que seja necessário recebê-lo imediatamente.

### Envio e recebimento

```go
ch <- "hello"
```

Envia um valor para o canal.

```go
message := <-ch
```

Recebe um valor do canal e o atribui à variável `message`.

### Exemplo: buffer de tamanho 1

```go
func main() {
    ch := make(chan string, 1)

    ch <- "A"
    ch <- "B"

    fmt.Println(<-ch)
}
```

Esse código entra em deadlock.

Por quê?

**Primeiro envio:**

```go
ch <- "A"
```

O buffer está vazio, então `"A"` é armazenado.

```text
Buffer: ["A"]
```

**Segundo envio:**

```go
ch <- "B"
```

O buffer já está cheio. Como ninguém recebeu `"A"`, o envio de `"B"` fica bloqueado.

**Impressão:**

```go
fmt.Println(<-ch)
```

Essa linha nunca é executada, porque a goroutine principal já está bloqueada no segundo envio.

O runtime detecta que não há goroutines capazes de prosseguir e reporta um deadlock.

### Como corrigir?

Uma possibilidade é receber o primeiro valor antes de enviar o segundo:

```go
func main() {
    ch := make(chan string, 1)

    ch <- "A"
    fmt.Println(<-ch)

    ch <- "B"
    fmt.Println(<-ch)
}
```

Saída:

```text
A
B
```

## 🧪 8.1. Praticando Channels com Buffer

Vamos entender como o tamanho do buffer influencia o envio e o recebimento de valores em um channel.

### Exemplo 1 — Liberando espaço no buffer

```go
ch := make(chan string, 1)

ch <- "hello"
fmt.Println(<-ch)

ch <- "hello 2"
fmt.Println(<-ch)
```

**Saída:**

```text
hello
hello 2
```

**O que acontece?**

1. Enviamos `"hello"` para o channel.
2. `<-ch` recebe e remove `"hello"` do buffer, liberando espaço.
3. Enviamos `"hello 2"` para o channel.
4. `<-ch` recebe e remove `"hello 2"`.

Como o buffer tem capacidade `1`, precisamos liberar espaço antes de enviar outro valor.

---

### Exemplo 2 — Buffer com capacidade para dois valores

```go
ch := make(chan string, 2)

ch <- "hello"
ch <- "hello 2"

fmt.Println(<-ch)
```

**Saída:**

```text
hello
```

**O que acontece?**

O buffer consegue armazenar dois valores, então os dois envios são concluídos sem precisar de um recebimento entre eles.

Quando executamos `<-ch`, recebemos o primeiro valor enviado: `"hello"`.

O valor `"hello 2"` continua armazenado no buffer.

```text
Buffer antes do recebimento:
┌─────────┬─────────────┐
│ hello   │ hello 2     │
└─────────┴─────────────┘

Buffer depois de receber "hello":
┌─────────┬─────────────┐
│ vazio   │ hello 2     │
└─────────┴─────────────┘
```

---

### Exemplo 3 — Recebendo todos os valores

```go
ch := make(chan string, 2)

ch <- "hello"
ch <- "hello 2"

fmt.Println(<-ch)
fmt.Println(<-ch)
```

**Saída:**

```text
hello
hello 2
```

Como o channel armazena os valores na ordem em que são enviados, eles também são recebidos nessa ordem: primeiro `"hello"`, depois `"hello 2"`.

Esse comportamento é conhecido como **FIFO (First In, First Out)** — primeiro a entrar, primeiro a sair.

---

### Exemplo 4 — Deadlock por excesso de envios

```go
ch := make(chan string, 1)

ch <- "hello"
ch <- "hello 2"

fmt.Println(<-ch)
```

**Resultado:** deadlock.

**Por quê?**

1. `"hello"` ocupa a única posição do buffer.
2. O envio de `"hello 2"` fica bloqueado porque o buffer está cheio.
3. O programa nunca chega ao `fmt.Println`, que poderia receber `"hello"` e liberar espaço.

O runtime detecta que a goroutine principal não consegue continuar e reporta um erro de deadlock.

---

### 📌 Regra geral

| Operação | Comportamento |
|---|---|
| `ch <- valor` | Envia um valor para o channel |
| `valor := <-ch` | Recebe e remove um valor do channel |
| `make(chan string, 1)` | Permite armazenar até um valor sem recebimento imediato |
| `make(chan string, 2)` | Permite armazenar até dois valores sem recebimento imediato |

**Importante:** o buffer determina quantos valores podem ficar armazenados sem que um recebimento seja necessário para liberar espaço. Quando o buffer está cheio, um novo envio fica bloqueado até que haja espaço disponível.

Um channel com buffer não exige uma goroutine por definição. Ele apenas permite que os envios e recebimentos sejam desacoplados até o limite de sua capacidade.

Outra possibilidade é utilizar uma goroutine para enviar os valores enquanto a goroutine principal os recebe.

**Importante:** channels não exigem sempre um buffer ou uma goroutine. Eles precisam de um fluxo de comunicação que possa progredir. Um channel sem buffer pode funcionar perfeitamente entre goroutines diferentes, ou mesmo dentro da mesma goroutine, desde que envio e recebimento possam se sincronizar.

---

## 🧠 9. Resumo rápido

| Conceito | Para que serve |
|---|---|
| `go.mod` | Declara o módulo e seus requisitos de dependências |
| `go.sum` | Registra checksums para verificar a integridade dos módulos |
| `go mod tidy` | Sincroniza os requisitos e os checksums do módulo |
| Struct | Agrupa campos relacionados |
| Ponteiro | Permite acessar um valor por seu endereço |
| Interface | Define métodos que um tipo precisa satisfazer |
| `context.Context` | Propaga cancelamento, prazos e valores associados à operação |
| `select` | Aguarda operações em canais |
| Goroutine | Executa uma função concorrentemente |
| Channel | Comunica e sincroniza goroutines |
| Channel com buffer | Armazena uma quantidade limitada de valores |
| `make` | Inicializa slices, maps e channels |

---

## 📌 Imports utilizados nos exemplos

Para executar os exemplos que envolvem esses conceitos, lembre-se de adicionar os imports necessários ao arquivo Go.

```go
import (
    "context"
    "fmt"
    "time"
)
```

Os exemplos de interfaces e structs não precisam necessariamente desses três imports; cada arquivo deve importar somente os pacotes que utiliza.

**Dica de estudo:** não tente decorar tudo de uma vez. Priorize compreender o comportamento de ponteiros, interfaces, contextos e channels escrevendo pequenos programas e observando o que acontece durante a execução.