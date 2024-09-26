package main

import (
	"programasGo/FPPDSemaforo"
	"fmt"
)

// variaveis
var readers int = 0;
var mutex *FPPDSemaforo.Semaphore = FPPDSemaforo.NewSemaphore(1)
var roomEmpty *FPPDSemaforo.Semaphore = FPPDSemaforo.NewSemaphore(1)

