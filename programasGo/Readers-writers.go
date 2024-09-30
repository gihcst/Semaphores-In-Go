package main

import (
	"fmt"
	"programasGo/FPPDSemaforo"
)

// Struct para Lightswitch
type Lightswitch struct {
	counter int
	mutex   *FPPDSemaforo.Semaphore
}

// Função para inicializar o Lightswitch
func NewLightswitch() *Lightswitch {
	return &Lightswitch{
		counter: 0,
		mutex:   FPPDSemaforo.NewSemaphore(1),
	}
}

// Função Lock para leitores
func (ls *Lightswitch) Lock(semaphore *FPPDSemaforo.Semaphore) {
	ls.mutex.Wait()
	ls.counter++
	if ls.counter == 1 {
		semaphore.Wait()
	}
	ls.mutex.Signal()
}

// Função Unlock para leitores
func (ls *Lightswitch) Unlock(semaphore *FPPDSemaforo.Semaphore) {
	ls.mutex.Wait()
	ls.counter--
	if ls.counter == 0 {
		semaphore.Signal()
	}
	ls.mutex.Signal()
}

// Variáveis
var roomEmpty *FPPDSemaforo.Semaphore = FPPDSemaforo.NewSemaphore(1)
var lightswitch = NewLightswitch() // Inicializando o lightswitch

// Leitores
func Readers(f chan struct{}, n int) {
	// Usando o lightswitch para controlar o acesso dos leitores
	lightswitch.Lock(roomEmpty)
	Read(n)
	fmt.Println(n, " Terminou de ler")
	lightswitch.Unlock(roomEmpty)
	f <- struct{}{}
}

// Escritores
func Writers(f chan struct{}, n int) {
	roomEmpty.Wait()
	Write(n)
	fmt.Println(n, " Terminou de escrever")
	roomEmpty.Signal()
	f <- struct{}{}
}

// Escrever
func Write(n int) {
	fmt.Println(n, " Escrevendo")
}

// Ler
func Read(n int) {
	fmt.Println(n, " Lendo")
}

func main() {
	fim := make(chan struct{})
	for i := 0; i < 10; i++ {
		go Readers(fim, i)
		go Writers(fim, i)
	}

	for i := 0; i < 20; i++ {
		<-fim
	}
}