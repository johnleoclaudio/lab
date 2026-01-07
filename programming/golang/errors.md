# Notes on errors

## Types 
1. Expected errors - occurs during normal operation. Must be handled gracefully.
  - database query timeout
  - network resource unavailable
  - bad user input
2. Unexpected errors - should NOT happen during normal operation. Most probably a bug/ logical error. Using panic is widely accepted. Using the language in an unintended way like:
  - accessing out-of-bounds index in a slice
  - trying to close an already-closed channel
But still recommended to return and gracefully handle unexpected errors in most cases. Execept when returning the error adds an unacceptable amount of error handling to the rest of the codebase.



## Ways to create an error 

```go
// 1
errOffset := 1
fmt.Errorf("failed to parse JSON (at character %d", errOffset)

// 2
errors.New("body contains badly-formed JSON")
```
