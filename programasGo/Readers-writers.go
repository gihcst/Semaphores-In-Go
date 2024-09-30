package main

import (
	"fmt"
	"programasGo/FPPDSemaforo"
)

type Lightswitch struct {
	counter int
	mutex   *FPPDSemaforo.Semaphore
}

func NewLightswitch() *Lightswitch {
	return &Lightswitch{
		counter: 0,
		mutex:   FPPDSemaforo.NewSemaphore(1),
	}
}

func (ls *Lightswitch) lock(semaphore *FPPDSemaforo.Semaphore) {
	ls.mutex.Wait()
	ls.counter++
	if ls.counter == 1 {
		semaphore.Wait()
	}
	ls.mutex.Signal()
}

func (ls *Lightswitch) unlock(semaphore *FPPDSemaforo.Semaphore) {
	ls.mutex.Wait()
	ls.counter--
	if ls.counter == 0 {
		semaphore.Signal()
	}
	ls.mutex.Signal()
}

var roomEmpty *FPPDSemaforo.Semaphore = FPPDSemaforo.NewSemaphore(1)
var readLightswitch = NewLightswitch() // usando lightswitch

func readers(f chan struct{}, n int) {
	readLightswitch.lock(roomEmpty)
	Read(n)
	fmt.Println(n, " terminou de ler!")
	readLightswitch.unlock(roomEmpty)
	f <- struct{}{}
}

func Writers(f chan struct{}, n int) {
	roomEmpty.Wait()
	Write(n)
	fmt.Println(n, " terminou de escrever!")
	roomEmpty.Signal()
	f <- struct{}{}
}

func Write(n int) {
	fmt.Println(n, " está escrevendo...")
}

func Read(n int) {
	fmt.Println(n, " está lendo...")
}

func main() {
	fim := make(chan struct{})
	for i := 0; i < 10; i++ {
		go readers(fim, i)
		go Writers(fim, i)
	}

	for i := 0; i < 20; i++ {
		<-fim
	}
}
