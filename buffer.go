package main

import (
	"bytes"
	"fmt"
)

func main() {
	s := bytes.NewBuffer([]byte("path/dir/file"))
	//s.Write([]byte("another"))
	//s.WriteTo(os.Stdout)
	//fmt.Println(string(s.Next(4)))
	//fmt.Println(string(s.Next(4)))
	word, _ := s.ReadBytes(byte('/'))
	fmt.Println(string(word))
	fmt.Println(string(s.String()))

}
