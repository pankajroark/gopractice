package main

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"hash/fnv"
	"os"
)

// Note that max size of string is 2bytes

type Node struct {
	offset uint32
	next   *Node
}

type Stringids struct {
	indexPath   string
	wal         *os.File
	walSize     uint32
	offsetTable []*Node
	capacity    uint32
	size        uint32
}

func (s *Stringids) init(path string) {
	s.indexPath = path
	s.capacity = 1024
	s.setup()
}

func (s *Stringids) setup() {
	s.wal, _ = os.OpenFile(s.indexPath, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0660)
	s.offsetTable = make([]*Node, s.capacity)
}

func (s *Stringids) loc(str string) uint32 {
	hash := fnv.New32()
	_, e := hash.Write([]byte(str))
	if e != nil {
		panic(e)
	}
	h := hash.Sum32()
	return h % uint32(len(s.offsetTable))
}

func (s *Stringids) add(str string) uint32 {
	slot := s.loc(str)
	fmt.Printf("slot %d\n", slot)
	offset := s.walSize
	binary.Write(s.wal, binary.LittleEndian, uint16(len(str)))
	n, err := s.wal.Write([]byte(str))
	if err != nil {
		panic(err)
	}
	s.walSize += uint32(n) + 2 // 2 bytes for size
	s.offsetTable[slot] = &Node{offset: offset, next: s.offsetTable[slot]}
	s.size += 1
	return offset
}

func (s *Stringids) strAtOffset(offset uint32) string {
	var size uint16
	ba := make([]byte, 2)
	s.wal.ReadAt(ba, int64(offset))
	reader := bytes.NewReader(ba)
	binary.Read(reader, binary.LittleEndian, &size)
	ba = make([]byte, size)
	offset += 2
	s.wal.ReadAt(ba, int64(offset))
	return string(ba)
}

func (s *Stringids) getId(str string) (uint32, error) {
	slot := s.loc(str)
	fmt.Printf("slot %d\n", slot)
	node := s.offsetTable[slot]
	for node != nil {
		tstr := s.strAtOffset(node.offset)
		if tstr == str {
			return node.offset, nil
		}
		node = node.next
	}
	return 0, errors.New("not found")
}

func (s *Stringids) clear() {
	err := os.Remove(s.indexPath)
	if err != nil {
		panic(err)
	}
	s.setup()
}

func main() {
	s := Stringids{}
	s.init("/tmp/testpath")
	s.clear()
	fmt.Printf("offset %d\n", s.add("some"))
	fmt.Printf("offset %d\n", s.add("other"))
	fmt.Printf("offset %d\n", s.add("thing"))
	fmt.Println(s.strAtOffset(0))
	fmt.Println(s.strAtOffset(6))
	fmt.Println(s.getId("other"))
	fmt.Println(s.getId("thing"))
	//fmt.Println
}
