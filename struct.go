package main

import "fmt"

type person struct {
	name string
	age  int
}

func main() {
	fmt.Println(person{"Bob", 20})
	fmt.Println(person{name: "Alice"})
	fmt.Println(person{name: "Fred", age: 30})
	fmt.Println(person{age: 35})

	s := person{name: "Sean", age: 50}
	fmt.Println(s.name)

	sp := &s
	fmt.Println(sp.age)

	sp.age = 51
	fmt.Println(sp.age)
	fmt.Println(s.age)

	cp := s
	cp.age = 52
	fmt.Println(cp.age)
	fmt.Println(s.age)
}
