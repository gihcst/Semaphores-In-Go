package main

import (
	"fmt"
	"programasGo/FPPDSemaforo"
	
)

// variaveis
var readers int = 0;
var mutex *FPPDSemaforo.Semaphore = FPPDSemaforo.NewSemaphore(1)
var roomEmpty *FPPDSemaforo.Semaphore = FPPDSemaforo.NewSemaphore(1)

//Leitores 
func Readers (f chan struct{}, n int){
	mutex.Wait()
	readers+=1
	if (readers ==1){
		roomEmpty.Wait()
	}
	mutex.Signal()
	Read(n)
	mutex.Wait()
	fmt.Println(n, " Terminou de ler")
	f <- struct{}{}
	readers -= 1
	if (readers == 0){
		roomEmpty.Signal()
	}
	mutex.Signal()
}
//Escritores
func Writers (f chan struct{}, n int ){
	roomEmpty.Wait()
	Write(n)
	fmt.Println(n, " Terminou de escrever")
	roomEmpty.Signal();
	f <- struct{}{}
}

//Escrever
func Write (n int){
	fmt.Println(n, " Escrevendo")
}
//Ler
func Read (n int){
	fmt.Println(n, " Lendo")
}

func main (){
	fim := make(chan struct{})
	for i := 0; i < 10; i++ {
		go Readers(fim, i)
		go Writers(fim, i)
	}

	for i := 0; i < 20; i++ {
		<-fim
	}
}
