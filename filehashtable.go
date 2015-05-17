package main

import "fmt"
import "bytes"

/*
File hash table stores data on disk. It is an immutable hashtable, key-values can only be added.
The design is based on open addressing hashtable with linear scan. Keys are assumed to be uint32.
Value is a byte array.
*/

// Let's start with a buffer based approach

type Fht struct {
	b    []byte
	size uint32
}

func CreateFHT(m map[uint32][]byte) {
	ht := Fht{}
	ht.b = make([]byte, 1024)
	for k, v := range m {
		ht.addKV(k, v)
	}
}

func (f *Fht) addKV(k uint, v []byte) {
	// calculate hash, check the slow, linear probe until found
}

func (f *Fht) get(key uint32) {
}
