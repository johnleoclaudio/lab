package main

import (
	"fmt"
	"time"
)

func doWork(id int) {
	fmt.Printf("Work %d started at %s\n", id, time.Now().Format("15:04:05"))
	time.Sleep(1 * time.Second)
	fmt.Printf("Work %d finished at %s\n", id, time.Now().Format("15:04:05"))
}

func main() {
	t := time.NewTimer(500 * time.Millisecond)
	for {
		select {
		case _, ok := <-t.C:
			if !ok {
				return
			}
			fmt.Println("time out!")
			return
		default:
			fmt.Println("waiting...")
		}
	}
}
