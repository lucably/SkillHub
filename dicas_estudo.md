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

# 10. HTTP Handlers, JSON e Testes de API no Go

## 10.1. O que é um Handler?

Um **handler** é uma função responsável por receber e tratar uma requisição HTTP, produzindo uma resposta para o cliente.

No Go, podemos criar handlers usando a biblioteca padrão `net/http`, sem precisar instalar um framework externo.

Exemplo:

```go
package main

import (
	"net/http"
)

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("API funcionando!"))
}

func main() {
	http.HandleFunc("/health", healthHandler)

	http.ListenAndServe(":8080", nil)
}
```

Nesse exemplo, quando alguém acessar `http://localhost:8080/health`, o Go executará a função `healthHandler`.

### Entendendo os parâmetros

**1. `w http.ResponseWriter`**

É o objeto utilizado para construir a resposta HTTP que será enviada ao cliente.

Por meio dele, podemos:

- Definir o status HTTP.
- Adicionar headers.
- Escrever o conteúdo da resposta.

Exemplo:

```go
w.WriteHeader(http.StatusOK)
w.Write([]byte("Sucesso"))
```

Nesse caso, o servidor responde com o status `200 OK` e o texto `Sucesso`.

**2. `r *http.Request`**

Contém as informações da requisição recebida pelo servidor, como:

- `r.Method`: método HTTP, como `GET`, `POST`, `PUT`, `PATCH` e `DELETE`.
- `r.URL`: informações sobre a URL acessada.
- `r.Header`: headers enviados pelo cliente.
- `r.Body`: corpo da requisição, utilizado, por exemplo, para receber um JSON.
- `r.PathValue("id")`: valor de um parâmetro de rota nomeado, quando configurado em um roteador compatível, como o `ServeMux` da biblioteca padrão nas versões recentes do Go.

Exemplo:

```go
func welcomeHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Bem-vindo!"))

	case http.MethodPost:
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte("Recurso criado!"))

	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}
```

O handler verifica o método HTTP e executa o comportamento correspondente.

O método `switch` é útil quando um mesmo endpoint aceita mais de um método HTTP. Em aplicações maiores, também podemos registrar handlers específicos para cada método e rota.

### Analogia com o frontend

No frontend, você pode ter uma função responsável por tratar o clique de um botão, validar os dados e chamar uma API.

No backend, o handler possui uma responsabilidade parecida: ele recebe a requisição HTTP, interpreta seus dados e encaminha o processamento para a camada adequada.

A diferença é que o handler trabalha no servidor e precisa produzir uma resposta HTTP.

---

## 10.2. Arquitetura de uma API: Router, Handler, Service e Repository

Uma API pode separar suas responsabilidades em diferentes camadas para facilitar a manutenção, os testes e a evolução do código.

### 1. Router

O router identifica qual rota foi acessada e encaminha a requisição para o handler correspondente.

Exemplo conceitual:

```text
GET /services       -> listarServices
POST /services      -> criarService
GET /services/{id}  -> buscarService
DELETE /services/{id} -> excluirService
```

Na biblioteca padrão do Go, podemos registrar rotas com `http.HandleFunc` ou utilizar `http.NewServeMux()`.

### 2. Handler

O handler é a camada que conversa diretamente com o protocolo HTTP.

Suas responsabilidades normalmente incluem:

- Ler parâmetros da URL.
- Verificar o método HTTP.
- Decodificar o JSON recebido.
- Validar a entrada em conjunto com as regras necessárias.
- Chamar o service.
- Converter o resultado em uma resposta HTTP.
- Definir o status e os headers da resposta.

O handler deve evitar concentrar todas as regras de negócio da aplicação.

### 3. Service

O service contém as regras de negócio.

Por exemplo, em um sistema de serviços profissionais, ele pode verificar:

- Se o preço informado é válido.
- Se o usuário tem autorização para realizar uma operação.
- Se um serviço pode ser excluído.
- Se determinada operação respeita as regras do negócio.

O service não deveria depender diretamente de detalhes HTTP, como `http.ResponseWriter`.

Isso facilita reutilizar a lógica em diferentes interfaces, como uma API HTTP, um worker ou uma tarefa agendada.

### 4. Repository

O repository é responsável pelo acesso e pela persistência dos dados.

Suas operações podem incluir:

- Buscar um registro no banco.
- Inserir um novo registro.
- Atualizar dados existentes.
- Excluir registros.

Ele pode utilizar SQL diretamente ou alguma biblioteca de acesso a banco de dados.

### Fluxo completo de uma requisição

```text
Cliente (Frontend, curl etc.)
          |
          v
        Router
          |
          v
        Handler
          |
          v
        Service
          |
          v
       Repository
          |
          v
       Banco de dados
```

A resposta percorre o caminho inverso: o resultado retorna ao handler, que monta a resposta HTTP e a envia ao cliente.

**Importante:** essa separação é uma convenção arquitetural, não uma exigência do Go. Para projetos pequenos, é possível começar com menos camadas e separá-las conforme a complexidade aumentar.

---

## 10.3. JSON no Go: Marshal, Unmarshal, Encode e Decode

Em APIs, frequentemente precisamos converter dados entre structs do Go e JSON.

O pacote utilizado para isso é:

```go
import "encoding/json"
```

Embora os quatro métodos trabalhem com JSON, eles são usados em situações diferentes.

### 10.3.1. Marshal vs. Unmarshal — trabalhando com `[]byte`

Essa dupla é utilizada quando o JSON está disponível como um conjunto de bytes na memória.

#### Marshal: Go → JSON

O método `json.Marshal()` converte um valor Go, como uma struct, para um JSON representado por `[]byte`.

Exemplo:

```go
type Pessoa struct {
	Nome  string `json:"nome"`
	Idade int    `json:"idade"`
}

func main() {
	p1 := Pessoa{
		Nome:  "Lucas",
		Idade: 30,
	}

	dadosBytes, err := json.Marshal(p1)
	if err != nil {
		panic(err)
	}

	fmt.Println(string(dadosBytes))
}
```

Saída:

```json
{"nome":"Lucas","idade":30}
```

Observe que o resultado de `json.Marshal()` é `[]byte`, e não uma string.

Por isso, utilizamos `string(dadosBytes)` para imprimir o conteúdo como texto.

O `Marshal` é útil quando precisamos:

- Gerar uma representação JSON na memória.
- Preparar dados para armazenamento ou transmissão.
- Produzir conteúdo JSON para logs ou outros processamentos.

#### Unmarshal: JSON → Go

O método `json.Unmarshal()` realiza o processo inverso: interpreta um JSON armazenado em `[]byte` e preenche uma variável Go.

Exemplo:

```go
dadosJSON := []byte(`{"nome":"Lucas","idade":30}`)

var p2 Pessoa

err := json.Unmarshal(dadosJSON, &p2)
if err != nil {
	panic(err)
}

fmt.Println(p2.Nome)
fmt.Println(p2.Idade)
```

Saída:

```text
Lucas
30
```

O segundo argumento é `&p2`, ou seja, o endereço da variável que será preenchida.

Isso é importante porque o `Unmarshal` precisa modificar a struct recebida.

**Resumo:**

| Método | Conversão | Entrada | Saída |
|---|---|---|---|
| `json.Marshal()` | Go → JSON | Struct ou outro valor Go | `[]byte` e `error` |
| `json.Unmarshal()` | JSON → Go | `[]byte` e endereço do destino | `error` |

Observação: strings JSON podem ser convertidas para `[]byte` antes de chamar `Unmarshal`, como no exemplo acima.

---

### 10.3.2. Encode vs. Decode — trabalhando com fluxos

`json.NewEncoder()` e `json.NewDecoder()` são úteis quando os dados são lidos ou escritos por meio de interfaces de fluxo, como `io.Reader` e `io.Writer`.

Um exemplo comum é trabalhar com requisições e respostas HTTP.

#### Encode: Go → JSON → fluxo de saída

O método `Encode()` transforma um valor Go em JSON e escreve o resultado no destino configurado no encoder.

Exemplo:

```go
func pessoaHandler(w http.ResponseWriter, r *http.Request) {
	pessoa := Pessoa{
		Nome:  "Lucas",
		Idade: 30,
	}

	w.Header().Set("Content-Type", "application/json")

	err := json.NewEncoder(w).Encode(pessoa)
	if err != nil {
		http.Error(
			w,
			"Erro ao gerar JSON",
			http.StatusInternalServerError,
		)
		return
	}
}
```

Nesse exemplo:

- `http.ResponseWriter` implementa `io.Writer`.
- `json.NewEncoder(w)` cria um encoder que escreve no destino `w`.
- `.Encode(pessoa)` converte a struct em JSON e escreve a resposta.

O cliente recebe um JSON semelhante a:

```json
{"nome":"Lucas","idade":30}
```

O `Encode()` também adiciona uma quebra de linha ao final do valor JSON.

É uma abordagem comum para responder a uma API sem precisar criar primeiro um `[]byte` com `json.Marshal()`.

#### Decode: fluxo de entrada → JSON → Go

O método `Decode()` lê o JSON de um fluxo e preenche a variável Go informada.

Exemplo:

```go
func criarPessoaHandler(w http.ResponseWriter, r *http.Request) {
	var pessoa Pessoa

	err := json.NewDecoder(r.Body).Decode(&pessoa)
	if err != nil {
		http.Error(
			w,
			"JSON inválido",
			http.StatusBadRequest,
		)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Pessoa recebida com sucesso"))
}
```

Nesse exemplo:

- `r.Body` contém o corpo da requisição HTTP.
- `json.NewDecoder(r.Body)` cria um decoder que lê desse fluxo.
- `.Decode(&pessoa)` interpreta o JSON e preenche a struct.
- Se o JSON estiver malformado ou não puder ser decodificado, o método retorna um erro.

Uma requisição POST poderia enviar:

```json
{
  "nome": "Lucas",
  "idade": 30
}
```

Depois do `Decode`, os valores estarão disponíveis em `pessoa.Nome` e `pessoa.Idade`.

**Atenção:** conseguir decodificar um JSON não significa que os dados são válidos para as regras do negócio. Por exemplo, uma idade negativa pode ser um JSON perfeitamente válido, mas a aplicação pode precisar rejeitá-la.

#### Resumo: Encode/Decode vs. Marshal/Unmarshal

| Método | Conversão | Origem ou destino típico |
|---|---|---|
| `Marshal` | Go → JSON | Retorna `[]byte` |
| `Unmarshal` | JSON → Go | Recebe `[]byte` |
| `Encode` | Go → JSON | Escreve em um `io.Writer` |
| `Decode` | JSON → Go | Lê de um `io.Reader` |

Uma regra prática:

- Se você já tem os bytes JSON na memória, considere `Marshal` ou `Unmarshal`.
- Se você está lendo de uma requisição HTTP ou escrevendo em uma resposta HTTP, considere `Decoder` ou `Encoder`.

Atenção à expressão “ler ou escrever aos poucos”: essas APIs trabalham com interfaces de fluxo e podem evitar que você precise carregar todo o JSON em um `[]byte` próprio. Isso não significa que todo valor JSON será necessariamente processado em pequenos blocos de forma incremental em todas as situações.

---

## 10.4. Testando uma API pelo terminal com `curl`

O `curl` permite enviar requisições HTTP manualmente. Ele é útil para testar uma API sem precisar criar uma interface frontend.

Antes de executar os comandos, certifique-se de que o servidor Go está em execução e escutando na porta `8080`.

Por exemplo:

```go
http.ListenAndServe(":8080", nil)
```

### 1. GET — consultar um endpoint

```bash
curl -i http://localhost:8080/welcome
```

Esse comando envia uma requisição GET para `/welcome`.

O método GET normalmente é utilizado para consultar informações.

A opção `-i` mostra os headers e o status HTTP da resposta, além do corpo retornado pelo servidor.

Uma resposta possível:

```http
HTTP/1.1 200 OK
Content-Type: text/plain; charset=utf-8

Bem-vindo!
```

O resultado exato depende do handler implementado.

### 2. POST — enviar uma requisição sem corpo

```bash
curl -i -X POST http://localhost:8080/welcome
```

A opção `-X POST` define o método HTTP como POST.

Esse comando não envia um corpo JSON. Ele apenas solicita ao endpoint que trate uma requisição POST.

O servidor pode aceitar a requisição, retornar `405 Method Not Allowed` ou responder de outra forma, conforme sua implementação.

### 3. POST — enviar dados JSON

No Windows, usando `curl.exe` no PowerShell, podemos executar:

```powershell
curl.exe -i -X POST http://localhost:8080/services -H "Content-Type: application/json" -d "{\"name\":\"Desenvolvimento de sites\",\"price\":1500}"
```

Esse comando envia um JSON para o endpoint `/services`.

O corpo enviado é:

```json
{
  "name": "Desenvolvimento de sites",
  "price": 1500
}
```

Entendendo as opções:

- `-i`: exibe os headers e o status HTTP da resposta.
- `-X POST`: define o método HTTP como POST.
- `-H "Content-Type: application/json"`: informa que o corpo da requisição está sendo enviado como JSON.
- `-d "..."`: envia os dados no corpo da requisição.

O backend pode receber esses dados usando `json.NewDecoder(r.Body).Decode(&service)`.

**Observação:** o comando acima está formatado para PowerShell do Windows. No Linux, macOS ou Git Bash, a forma de escapar as aspas pode ser diferente.

### 4. POST — enviar um JSON inválido

```powershell
curl.exe -i -X POST http://localhost:8080/services -H "Content-Type: application/json" -d "{\"name\":\"Desenvolvimento de sites\",\"price\":}"
```

Nesse exemplo, o campo `price` está sem um valor válido.

O corpo contém um JSON malformado:

```json
{"name":"Desenvolvimento de sites","price":}
```

O decoder deve retornar um erro ao tentar interpretar esse conteúdo.

Uma API bem implementada deve tratar esse erro e responder com um status apropriado, normalmente `400 Bad Request`.

O comando serve para verificar se o backend trata entradas inválidas corretamente, em vez de falhar inesperadamente.

### 5. DELETE — excluir um recurso

```powershell
curl.exe -i -X DELETE http://localhost:8080/services/123
```

Esse comando envia uma requisição DELETE para `/services/123`.

A intenção é excluir o recurso cujo identificador é `123`.

Em uma implementação real, o backend deve:

1. Extrair o identificador da rota.
2. Validar o identificador.
3. Verificar as permissões necessárias.
4. Executar a regra de negócio.
5. Excluir o registro, caso a operação seja permitida.
6. Retornar um status HTTP apropriado.

O status pode ser `204 No Content` quando a exclusão é concluída sem corpo de resposta, ou outro status conforme o contrato da API.

O simples envio do comando não garante que o recurso será excluído. Isso depende do handler e das camadas chamadas por ele.

### 6. PATCH — atualizar parcialmente um recurso

```powershell
curl.exe -i -X PATCH http://localhost:8080/services/123 -H "Content-Type: application/json" -d "{\"price\":1800}"
```

Esse comando envia uma requisição PATCH para atualizar o preço do serviço `123`.

O JSON enviado é:

```json
{
  "price": 1800
}
```

PATCH costuma ser utilizado para atualizações parciais. Nesse caso, a intenção é alterar o preço sem necessariamente enviar todos os campos do serviço.

Uma implementação pode utilizar um DTO específico para atualização parcial:

```go
type UpdateServiceDTO struct {
	Price *float64 `json:"price"`
}
```

O ponteiro permite diferenciar alguns casos:

- `Price == nil`: o campo não foi informado, ou foi enviado como `null`.
- `Price != nil`: o campo contém um valor numérico, inclusive `0`.

Se for necessário distinguir um campo ausente de um campo explicitamente enviado como `null`, será preciso adotar uma representação ou lógica adicional.

Após decodificar o JSON, o handler deve encaminhar a atualização ao service, que aplicará as regras de negócio antes de persistir as alterações.

### 7. PUT — atualizar um recurso completo

```powershell
curl.exe -i -X PUT http://localhost:8080/services/123 -H "Content-Type: application/json" -d "{\"name\":\"Desenvolvimento de sistemas\",\"price\":2000}"
```

Esse comando envia uma requisição PUT para atualizar o serviço de identificador `123`.

O JSON enviado é:

```json
{
  "name": "Desenvolvimento de sistemas",
  "price": 2000
}
```

O método **PUT** normalmente é utilizado para substituir a representação de um recurso por uma nova versão. Por isso, é comum enviar todos os campos necessários para representar o recurso atualizado.

Por exemplo, se um serviço possui os campos `name` e `price`, o PUT pode receber ambos os campos para substituir os valores anteriores.

Em uma implementação típica, o backend deve:

1. Extrair o identificador da URL.
2. Decodificar o JSON recebido usando `json.NewDecoder(r.Body).Decode(&service)`.
3. Validar os dados enviados.
4. Verificar se o recurso existe e se o usuário tem permissão para alterá-lo.
5. Executar as regras de negócio no service.
6. Persistir a atualização no banco de dados.
7. Retornar o status HTTP adequado.

Uma resposta possível é `200 OK`, caso a API retorne o recurso atualizado, ou `204 No Content`, caso conclua a operação sem retornar um corpo.

**Diferença entre PUT e PATCH:**

- **PUT:** normalmente substitui a representação do recurso, recebendo os campos necessários para a nova versão.
- **PATCH:** normalmente altera apenas os campos especificados na requisição.

Por exemplo, se o serviço possui nome e preço:

- PUT enviando `name` e `price`: representa o serviço completo após a atualização.
- PATCH enviando somente `price`: solicita a alteração apenas do preço.

A semântica exata depende do contrato da API. Em geral, PUT também é considerado idempotente: repetir a mesma requisição deve produzir o mesmo estado final do recurso, sem acumular alterações a cada repetição.

Importante: enviar o comando `curl` não implementa a atualização automaticamente. O backend precisa ter um handler que aceite PUT e execute a operação.


---

## 10.5. Boas práticas ao criar handlers

Ao desenvolver uma API em Go, algumas práticas ajudam a manter o código seguro e organizado:

1. **Trate erros:** verifique o erro retornado por `json.NewDecoder(r.Body).Decode()`, `json.Marshal()` e outras operações que possam falhar.
2. **Use os status HTTP adequados:** por exemplo, `200 OK`, `201 Created`, `204 No Content`, `400 Bad Request`, `404 Not Found` e `500 Internal Server Error`.
3. **Separe as responsabilidades:** deixe o handler cuidar do HTTP e encaminhe as regras de negócio ao service.
4. **Valide os dados:** JSON válido não significa dados válidos para o negócio.
5. **Não exponha erros internos:** evite retornar detalhes sensíveis do banco de dados ou da infraestrutura ao cliente.
6. **Cuidado ao escrever a resposta:** depois que o status ou o corpo HTTP já foram enviados, não é possível simplesmente substituí-los por outra resposta.
7. **Limite o tamanho do corpo:** em endpoints que recebem JSON, considere limitar o tamanho de `r.Body` para evitar requisições excessivamente grandes.
8. **Teste cenários inválidos:** envie JSON malformado, campos ausentes, valores fora do esperado e métodos HTTP não permitidos.

---

## 10.6. Resumo geral

| Conceito | Responsabilidade |
|---|---|
| Router | Encaminha a requisição para a rota adequada |
| Handler | Interpreta HTTP e monta a resposta |
| Service | Executa as regras de negócio |
| Repository | Acessa e persiste os dados |
| `json.Marshal()` | Converte um valor Go em `[]byte` JSON |
| `json.Unmarshal()` | Converte `[]byte` JSON em um valor Go |
| `json.NewEncoder()` | Cria um encoder para escrever JSON em um `io.Writer` |
| `json.NewDecoder()` | Cria um decoder para ler JSON de um `io.Reader` |
| `curl` | Permite testar endpoints HTTP pelo terminal |
| `GET` | Consultar recursos |
| `POST` | Criar recursos ou executar operações |
| `PATCH` | Atualizar parcialmente um recurso |
| `PUT` | Substituir a representação de um recurso |
| `DELETE` | Solicitar a exclusão de um recurso |

**O principal aprendizado:** uma API Go recebe uma requisição HTTP, encaminha o processamento para o handler, converte os dados quando necessário, executa as regras de negócio nas camadas adequadas e retorna uma resposta HTTP com o status e o conteúdo correspondentes.
