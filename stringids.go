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

type OffsetTable struct {
	offsetTable []*Node
	size        uint32
}

func NewOffsetTable() *OffsetTable {
	ot := make([]*Node, 1024)
	o := OffsetTable{offsetTable: ot, size: 0}
	return &o
}

func (t *OffsetTable) put(hash, offset uint32) {
	slot := hash % uint32(len(t.offsetTable))
	t.offsetTable[slot] = &Node{offset: offset, next: t.offsetTable[slot]}
	t.size += 1
}

func (t *OffsetTable) get(hash uint32) *Node {
	slot := hash % uint32(len(t.offsetTable))
	fmt.Printf("slot %d\n", slot)
	return t.offsetTable[slot]
}

type Stringids struct {
	indexPath   string
	wal         *os.File
	walSize     uint32
	offsetTable *OffsetTable
}

func NewStringids(path string) *Stringids {
	wal, _ := os.OpenFile(path, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0660)
	fi, e := wal.Stat()
	if e != nil {
		panic(e)
	}
	walSize := uint32(fi.Size())
	offsetTable := NewOffsetTable()
	return &Stringids{indexPath: path, wal: wal, walSize: walSize, offsetTable: offsetTable}

}

func (s *Stringids) reset() {
	s.wal, _ = os.OpenFile(s.indexPath, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0660)
	fi, e := s.wal.Stat()
	if e != nil {
		panic(e)
	}
	s.walSize = uint32(fi.Size())
	s.offsetTable = NewOffsetTable()
}

func (s *Stringids) hash(str string) uint32 {
	hash := fnv.New32()
	_, e := hash.Write([]byte(str))
	if e != nil {
		panic(e)
	}
	return hash.Sum32()
}

func (s *Stringids) writeToWal(str string) uint32 {
	offset := s.walSize
	binary.Write(s.wal, binary.LittleEndian, uint16(len(str)))
	n, err := s.wal.Write([]byte(str))
	if err != nil {
		panic(err)
	}
	s.walSize += uint32(n) + 2 // 2 bytes for size
	return offset
}

func (s *Stringids) add(str string) uint32 {
	offset := s.writeToWal(str)
	s.offsetTable.put(s.hash(str), offset)
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
	node := s.offsetTable.get(s.hash(str))
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
	s.reset()
}

func main() {
	s := NewStringids("/tmp/testpath")
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
