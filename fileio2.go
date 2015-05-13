package main

import (
	"fmt"
	"os"
)

func main() {
	var err error
	var f *os.File
	f, err = os.OpenFile("/tmp/goiotest", os.O_RDWR|os.O_APPEND, 0660)
	if err != nil {
		panic(err)
	}

	var n int
	n, err = f.Write([]byte("some"))
	if err != nil {
		panic(err)
	}
	fmt.Println(n)
	n, err = f.Write([]byte("other"))
	if err != nil {
		panic(err)
	}
	fmt.Println(n)

	b := make([]byte, 5)
	n, err = f.ReadAt(b, 4)
	fmt.Println(string(b))

	// try to read other
}
