# Go's concurrency

- There are two ways to achieve concurrent programming in Go: Primitives and Channels

## Primitives - Protecting shared memory
- used when your goroutines needs to share memory safely
- Mutexes, WaitGroups, Atomic operations, and Cond variables


## Channels - Transferring ownership
- when you need to pass data between goroutines
- "Don't communicate by sharing memory; Share memory by communicating"
