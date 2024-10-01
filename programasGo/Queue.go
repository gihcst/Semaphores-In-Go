package main

import (
	"fmt"
	"programasGo/FPPDSemaforo"
	"sync"
)

var liders = 0
var seguidores = 0
var mutex = FPPDSemaforo.NewSemaphore(1) // garante q a modificação de liders e seguidores seja acessada por uma goroutine por vez
// filas de espera q impedem que um lider ou seguidor dancem sozinhos
var liderQueue = FPPDSemaforo.NewSemaphore(0)
var seguidorQueue = FPPDSemaforo.NewSemaphore(0)
var rendezvous = FPPDSemaforo.NewSemaphore(0) //redenzvous que sincroniza o termino da dança, garantindo que deixaram de dançar ao mesmo tempo
var wg sync.WaitGroup                         //garante que todas as goroutines executem antes de terminar

func lider() {
	defer wg.Done()

	mutex.Wait() //acessa as variaveis para ver se há seguidores disponíveis e bloqueia o acesso de outras goroutines
	if seguidores > 0 {
		seguidores--           //decrementa um seguidor da fila de espera
		seguidorQueue.Signal() //acorda o seguidor pra formar um par com algum lider
	} else {
		liders++ //se n houver seguidores, então o lider entra em espera
		mutex.Signal()
		liderQueue.Wait() //espera ate que um seguidor esteja livre para dançar!
	}
	dance("lider") //simula a dança
	rendezvous.Wait()
	mutex.Signal() //libera o acesso
}

func seguidor() {
	defer wg.Done()

	mutex.Wait()
	if liders > 0 {
		liders--
		liderQueue.Signal() //acorda um lider pra formar par
	} else {
		seguidores++ //entra na fila se nao houver lideres disponiveis
		mutex.Signal()
		seguidorQueue.Wait() //espera por um lider
	}
	dance("seguidor")   //simula a dança com um lider
	rendezvous.Signal() // sinal de q a dança foi concluida
}

func dance(role string) {
	fmt.Println(role, "está dançando no momento!!") //imprime quem esta dançando no momento
}

func main() {
	wg.Add(8) //adc lideres e seguidores ao waitgroup

	for i := 0; i < 4; i++ {
		go lider()
		go seguidor()
	}
	wg.Wait()
}
