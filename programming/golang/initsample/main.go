package main

import "fmt"

func init() {
	fmt.Println("init")
}

func main() {
	fmt.Println("main")
}

var a = func() int {
	fmt.Println("var")
	return 0
}()
