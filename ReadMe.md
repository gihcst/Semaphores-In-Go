# Concurrency Algorithms with Semaphores in Go

This repository contains multiple concurrent programs implemented in Go, all making use of semaphores for synchronization. Below is an explanation of each program along with its respective functionality.

## 1. Dining Savages Problem

This program simulates the Dining Savages problem, where a group of savages shares a pot of food and a cook prepares more food whenever the pot is empty. The savages eat the food from the pot, and when it's empty, they notify the cook to refill it.

### Logic:
- **Savages**: Take portions of food from the pot. If the pot is empty, they notify the cook.
- **Cook**: Refills the pot with `M` portions of food whenever it's notified that the pot is empty.

### Semaphores:
- **Mutex**: Ensures that only one savage accesses the pot at a time.
- **Empty Pot**: Signals the cook to refill the pot when it's empty.
- **Full Pot**: Signals the savages that the pot has been refilled.

### Output:
The program outputs which savage is eating and when the cook is refilling the pot.

## 2. Dance Problem (Leaders and Followers)

This program simulates a ballroom dance where leaders and followers form pairs to dance. A leader can only dance with a follower and vice versa.

### Logic:
- **Leaders**: Wait for followers to be available to form a pair.
- **Followers**: Wait for leaders to be available to form a pair.
- Both must dance together, and the dance ends synchronously.

### Semaphores:
- **Mutex**: Controls access to the count of leaders and followers.
- **LeaderQueue**: Keeps leaders waiting when no followers are available.
- **FollowerQueue**: Keeps followers waiting when no leaders are available.
- **Rendezvous**: Synchronizes the end of the dance so that both leader and follower stop dancing at the same time.

### Output:
The program shows which leader and follower are dancing together.

## 3. Readers-Writers Problem

This program simulates the Readers-Writers problem, where readers and writers compete for access to a shared resource (e.g., a database). The goal is to allow multiple readers to read simultaneously, but writers require exclusive access.

### Logic:
- **Readers**: Multiple readers can read the resource simultaneously.
- **Writers**: Only one writer can write at a time, and no readers can read while writing.

### Semaphores:
- **RoomEmpty**: Ensures that the room is empty for a writer to write.
- **ReadLightswitch**: Implements the lightswitch pattern, allowing multiple readers to enter or exit the room.

### Output:
The program shows when a reader is reading and when a writer is writing.

## 4. Santa Claus Problem

This program simulates the Santa Claus problem, where Santa must help elves or deliver presents with his reindeer. Santa only works when he is either called by elves or when all the reindeer are ready.

### Logic:
- **Santa**: Helps elves or prepares the sleigh for reindeer when necessary.
- **Elves**: Wait for Santa's help when three elves are ready.
- **Reindeers**: Wait for Santa to prepare the sleigh when all nine reindeer have arrived.

### Semaphores:
- **Mutex**: Controls access to the counts of elves and reindeers.
- **SantaSem**: Wakes Santa when elves or reindeers need help.
- **ReindeerSem**: Ensures that reindeers wait until Santa is ready to prepare the sleigh.
- **ElfTex**: Ensures that only three elves ask for help from Santa at a time.

### Output:
The program shows when Santa is helping elves or preparing the sleigh, as well as when elves and reindeer are ready.

## Running the Programs

To run any of these programs, use the `go run` command:

```bash
go run <filename.go>
```

Ensure you have installed the `FPPDSemaforo` package in your Go environment, as it is required for semaphore operations.