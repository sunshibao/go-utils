package main

import "fmt"

func add(x, y int) int {
	res := 0
	res = x + y
	return res
}

func main() {
	fmt.Println(add(1, 2))
}
