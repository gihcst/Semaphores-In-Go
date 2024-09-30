package main

import (
	"fmt"
	"programasGo/FPPDSemaforo"
)

var rounds int = 3 // n de vezes que o cozinheiro enche o pote (pode ser qualquer n)
var M int = 10     // n de porções que o cozinheiro coloca no pote (pode ser qualquer n)
var servings int = 0
var mutex = FPPDSemaforo.NewSemaphore(1) // mutex para acessar o pote
var emptyPot = FPPDSemaforo.NewSemaphore(0)
var fullPot = FPPDSemaforo.NewSemaphore(0)

func cook(f chan struct{}) {
	for r := 0; r < rounds; r++ { // coloquei rodadas pra não ficar infinito
		emptyPot.Wait()
		putServingsInPot(M)
		servings = M     // atualiza n de porções
		fullPot.Signal() // avisa que encheu o pote
	}
	f <- struct{}{}
}

func savages(f chan struct{}, n int) {
	for r := 0; r < rounds; r++ { // coloquei rodadas pra não ficar infinito
		mutex.Wait()
		if servings == 0 {
			emptyPot.Signal() // se pegou a ultima porção, avisa o cozinheiro
			fullPot.Wait()
			servings = M // atualiza n de porções
		}
		servings -= 1 // pega uma porção
		getServingFromPot(n)
		mutex.Signal()

		eat(n)
	}
	f <- struct{}{}
}

func putServingsInPot(m int) {
	fmt.Printf(">>> cozinheiro encheu o pote de novo (%d porções) <<<\n", m)
}

func getServingFromPot(n int) {
	fmt.Println(n, "pegou uma porção do pote...")
}

func eat(n int) {
	fmt.Println(n, "comeu!")
}

func main() {
	fim := make(chan struct{})

	go cook(fim)

	for i := 0; i < 10; i++ {
		go savages(fim, i)
	}

	for i := 0; i < 11; i++ {
		<-fim
	}
}
