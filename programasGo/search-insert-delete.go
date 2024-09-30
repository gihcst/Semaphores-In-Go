package main

import (
	"fmt"
	"sync"
)

// Define the Lightswitch structure for managing multiple threads
type Lightswitch struct {
	counter int
	mutex   sync.Mutex
}

// Wait increments the counter and locks the semaphore if this is the first searcher/inserter
func (ls *Lightswitch) Wait(sem *sync.Mutex) {
	ls.mutex.Lock()
	defer ls.mutex.Unlock()
	ls.counter++
	if ls.counter == 1 {
		sem.Lock() // First searcher/inserter locks the semaphore
	}
}

// Signal decrements the counter and unlocks the semaphore if this is the last searcher/inserter
func (ls *Lightswitch) Signal(sem *sync.Mutex) {
	ls.mutex.Lock()
	defer ls.mutex.Unlock()
	ls.counter--
	if ls.counter == 0 {
		sem.Unlock() // Last searcher/inserter unlocks the semaphore
	}
}

// Define a Node for the singly linked list
type Node struct {
	value int
	next  *Node
}

// Define the List structure
type List struct {
	head  *Node
	mutex sync.Mutex
}

// Add adds a new value to the end of the list
func (l *List) Add(value int) {
	l.mutex.Lock()
	defer l.mutex.Unlock()

	newNode := &Node{value: value}
	if l.head == nil {
		l.head = newNode
	} else {
		curr := l.head
		for curr.next != nil {
			curr = curr.next
		}
		curr.next = newNode
	}
}

// Delete deletes the first occurrence of a value in the list
func (l *List) Delete(value int) bool {
	l.mutex.Lock()
	defer l.mutex.Unlock()

	if l.head == nil {
		return false
	}
	if l.head.value == value {
		l.head = l.head.next
		return true
	}

	curr := l.head
	for curr.next != nil {
		if curr.next.value == value {
			curr.next = curr.next.next
			return true
		}
		curr = curr.next
	}
	return false
}

// Print prints the list
func (l *List) Print() {
	l.mutex.Lock()
	defer l.mutex.Unlock()

	curr := l.head
	for curr != nil {
		fmt.Print(curr.value, " -> ")
		curr = curr.next
	}
	fmt.Println("nil")
}

// Define semaphores and switches
var insertMutex sync.Mutex
var noSearcher sync.Mutex
var noInserter sync.Mutex
var searchSwitch Lightswitch
var insertSwitch Lightswitch

// Searcher function
func searcher(list *List, value int, wg *sync.WaitGroup) {
	defer wg.Done()

	searchSwitch.Wait(&noSearcher) // Ensure no deleter is active
	defer searchSwitch.Signal(&noSearcher)

	// Critical section: search the list
	list.mutex.Lock()
	curr := list.head
	found := false
	for curr != nil {
		if curr.value == value {
			found = true
			break
		}
		curr = curr.next
	}
	list.mutex.Unlock()

	if found {
		fmt.Printf("Searcher found %d in the list.\n", value)
	} else {
		fmt.Printf("Searcher did not find %d in the list.\n", value)
	}
}

// Inserter function
func inserter(list *List, value int, wg *sync.WaitGroup) {
	defer wg.Done()

	insertSwitch.Wait(&noInserter) // Ensure no deleter is active
	insertMutex.Lock()             // Ensure only one inserter at a time
	defer insertMutex.Unlock()
	defer insertSwitch.Signal(&noInserter)

	// Critical section: insert into the list
	list.Add(value)
	fmt.Printf("Inserter added %d to the list.\n", value)
}

// Deleter function
func deleter(list *List, value int, wg *sync.WaitGroup) {
	defer wg.Done()

	noSearcher.Lock() // Ensure no searcher is active
	noInserter.Lock() // Ensure no inserter is active
	defer noSearcher.Unlock()
	defer noInserter.Unlock()

	// Critical section: delete from the list
	if list.Delete(value) {
		fmt.Printf("Deleter removed %d from the list.\n", value)
	} else {
		fmt.Printf("Deleter did not find %d in the list.\n", value)
	}
}

func main() {
	var wg sync.WaitGroup

	// Initialize the list
	list := &List{}

	// Launching some searchers, inserters, and deleters
	wg.Add(5)
	go inserter(list, 1, &wg)
	go inserter(list, 2, &wg)
	go searcher(list, 1, &wg)
	go deleter(list, 1, &wg)
	go searcher(list, 3, &wg)

	wg.Wait()

	// Final list state
	fmt.Println("Final list:")
	list.Print()
}
