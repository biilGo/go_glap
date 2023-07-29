package main

import "fmt"

func NewTwiceFuncClosurce(x int) func() int {
	return func() int {
		x *= 2
		return x
	}
}

func main() {
	fnTwice := NewTwiceFuncClosurce(1)

	fmt.Println(fnTwice())
	fmt.Println(fnTwice())
	fmt.Println(fnTwice())
}
