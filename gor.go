package main

import "fmt"

func f(from string) {
	for i := 0; i < 30000; i++ {
		fmt.Println(from, ":", i)
	}
}

func main() {
	go f("goroutine")
	go f("r2")
	go f("r3")

	go func(msg string) {
		fmt.Println(msg)
	}("going")

	var input string
	fmt.Scanln(&input)
	fmt.Println("done")
}
