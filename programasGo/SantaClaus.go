package main

import (
	"fmt"
	"programasGo/FPPDSemaforo"
)

var rounds int = 3
var elves int = 0
var reindeers int = 0
var santaSem = FPPDSemaforo.NewSemaphore(0)
var reindeerSem = FPPDSemaforo.NewSemaphore(0)
var elfTex = FPPDSemaforo.NewSemaphore(1) // mutex para elfos
var mutex = FPPDSemaforo.NewSemaphore(1)

func santa(f chan struct{}) {
	for r := 0; r < rounds; r++ { // coloquei rodadas pra não ficar infinito
		santaSem.Wait()
		mutex.Wait()
		if reindeers == 9 { // se todas as renas chegaram, prepara o trenó
			prepareSleight()
			for i := 0; i < 9; i++ {
				reindeerSem.Signal()
			}
		} else if elves == 3 { // chama o papai noel para ajudar caso tenha três elfos
			helpElves()
		}
		mutex.Signal()
	}
	f <- struct{}{}
}

func reindeer(f chan struct{}, n int) {
	mutex.Wait()
	reindeers += 1
	if reindeers == 9 {
		santaSem.Signal()
	}
	mutex.Signal()

	// as renas ficam esperando o papai noel preparar o trenó
	reindeerSem.Wait()
	getHitched(n)
	f <- struct{}{}
}

func elf(f chan struct{}, n int) {
	for r := 0; r < rounds; r++ { // coloquei rodadas pra não ficar infinito
		elfTex.Wait()
		mutex.Wait()
		elves += 1
		fmt.Println(">>> elfo entrou <<<")
		if elves == 3 {
			santaSem.Signal()
		} else {
			elfTex.Signal()
		}
		mutex.Signal()

		getHelp(n)

		mutex.Wait()
		elves -= 1
		fmt.Println(">>> elfo saiu <<<")
		if elves == 0 {
			elfTex.Signal()
		}
		mutex.Signal()
	}
	f <- struct{}{}
}

func prepareSleight() {
	fmt.Println("papai noel está preparando o trenó!")
}

func getHitched(n int) {
	fmt.Println(n, "rena está pronta para partir!")
}

func helpElves() {
	fmt.Println("papai noel está ajudando os elfos!")
}

func getHelp(n int) {
	fmt.Println(n, "elfo está pedindoo ajuda!")
}

func main() {
	fim := make(chan struct{})

	go santa(fim)

	for i := 0; i < 9; i++ {
		go reindeer(fim, i)
	}

	for i := 0; i < 10; i++ {
		go elf(fim, i)
	}

	for i := 0; i < 20; i++ {
		<-fim
	}
}
