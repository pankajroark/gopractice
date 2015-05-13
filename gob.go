package main

import (
	"bytes"
	"encoding/gob"
	"fmt"
)

func main() {
	// map[string][]string
	b := new(bytes.Buffer)
	e := gob.NewEncoder(b)

	m := make(map[string][]string)
	m["some"] = []string{"one", "two"}

	err := e.Encode(m)
	if err != nil {
		panic(err)
	}

	var decodeMap map[string][]string
	d := gob.NewDecoder(b)
	d.Decode(&decodeMap)

	fmt.Printf("%v\n", decodeMap)
}
