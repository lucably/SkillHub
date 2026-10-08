package main

import (
	"context"
	"fmt"
	"time"
)

// type PaymentProcessor interface {
// 	Pay(amount float64) error
// }

// type PixPayment struct{}
// type CreditCardPayment struct{}

// func (p PixPayment) Pay(amount float64) error {
// 	fmt.Println("Processing PIX payment:", amount)
// 	return nil
// }

// func (c CreditCardPayment) Pay(amount float64) error {
// 	fmt.Println("Processing Credit Card payment:", amount)
// 	return nil
// }

// func processPayment(payment PaymentProcessor, amount float64) error {
// 	return payment.Pay(amount)
// }

// type Product struct {
// 	Name  string
// 	Price float64
// 	Stock int
// }

// func (p Product) IsAvailable() bool {
// 	return p.Stock >= 1
// }

// func (p *Product) AddStock(quantity int) {
// 	p.Stock += quantity
// }

// func withdraw(stock int, quantity int) (int, error) {
// 	if quantity <= 0 {
// 		return 0, errors.New("invalid quantity")
// 	}

// 	if quantity > stock {
// 		return 0, errors.New("insufficient stock")
// 	}

// 	return stock - quantity, nil
// }

func doSomething(ctx context.Context) error {
	select {
	case <-ctx.Done():
		fmt.Println("Operação foi cancelada.")
		return ctx.Err()
	case <-time.After(5 * time.Second):
		fmt.Println("Operação Terminou.")
		return ctx.Err()
	}
}

func worker(ch chan string) {
	time.Sleep(2 * time.Second)

	ch <- "Process finished"
}

func main() {

	// product := Product{
	// 	Name:  "Mechanical Keyboard",
	// 	Price: 350,
	// 	Stock: 10,
	// }

	// fmt.Println("IsAvailable: ", product.IsAvailable())
	// product.AddStock(15)
	// fmt.Println("New Stock: ", product.Stock)

	// value, err := withdraw(1, 3)
	// if err != nil {
	// 	fmt.Println("Error: ", err)
	// 	return
	// }
	// fmt.Println("value: ", value)

	// processPayment(PixPayment{}, 100)
	// processPayment(CreditCardPayment{}, 200)

	// ctx, cancel := context.WithTimeout(
	// 	context.Background(),
	// 	2*time.Second,
	// )

	// defer cancel()

	// doSomething(ctx)

	/* EXERCICIO CHANNEL e GOROUTINES
	ch := make(chan string, 1)
	//Aqui utiliza o goroutine falando para executar uma função de forma concorrente, ou seja, continua a func main juntamente com o worker.
	go worker(ch)

	//Ja aqui seria um await. Ou seja, estamos falando para ESPERAR (travar, mas outras goroutines podem continuar executando) aqui até que o ch tenha valor "Espere até que alguém coloque um valor nesse channel."
	message := <-ch
	fmt.Println(message)
	*/
}
