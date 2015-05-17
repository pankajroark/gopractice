package ds

// todo use tabulation hashing
import (
	"encoding/binary"
	"errors"
	"fmt"
	"hash/fnv"
)

/*
 Int hash table is an uint32 -> uint32 hashtable based on open addressing and
 linear probing. The buffer of the table grows as needed.
*/

const SlotSize = 9 // flag + key + value

// If we scan more than this fraction of buffer while trying to find a slot
// we should resize.
const ScanFactor = 10

const EmptySlot = 0x00
const FullSlot = 0x01

type IntHashTable struct {
	ba []byte
}

// @param initialCapacity initial capacity for storing number of key value pairs
func CreateIntHashTable(initialCapacity uint32) IntHashTable {
	iht := IntHashTable{}
	iht.ba = make([]byte, initialCapacity*SlotSize)
	return iht
}

func hash(k uint32) uint32 {
	// todo: use tabulation hashing
	h := fnv.New32()
	e := binary.Write(h, binary.LittleEndian, k)
	if e != nil {
		panic(e)
	}
	return h.Sum32()
}

// Simply write th 9 bytes to byte buffer starting at slot location
func writeKV(buffer []byte, k, v, loc uint32) {
	buffer[loc] = FullSlot
	putUInt32(buffer, loc+1, k)
	putUInt32(buffer, loc+5, v)
}

func put(buffer []byte, k, v uint32) error {
	numSlots := uint32(cap(buffer) / SlotSize)
	h := hash(k)
	slot := h % numSlots
	slotByteLoc := slot * SlotSize
	firstSlot := slot
	for buffer[slotByteLoc] != EmptySlot {
		// explicit check in case we add another state for slot in future
		if buffer[slotByteLoc] == FullSlot {
			readKey := readUInt32(buffer, slotByteLoc+1)
			if readKey == k {
				// key already exists just overwrite it
				break
			}
		}
		slot += 1
		// slot is monotonically increasing
		if slot-firstSlot > numSlots/ScanFactor {
			return errors.New("Buffer too full.")
		}
		// slotByteLoc wraps around
		slotByteLoc = (slot % numSlots) * SlotSize
	}
	writeKV(buffer, k, v, slotByteLoc)
	return nil
}

// LittleEndian
func putUInt32(ba []byte, offset, i uint32) {
	ba[offset] = byte(i)
	ba[offset+1] = byte(i >> 8)
	ba[offset+2] = byte(i >> 16)
	ba[offset+3] = byte(i >> 24)
}

// LittleEndian
func readUInt32(ba []byte, offset uint32) uint32 {
	return uint32(ba[offset]) |
		uint32(ba[offset+1])<<8 |
		uint32(ba[offset+2])<<16 |
		uint32(ba[offset+3])<<24
}

func (iht *IntHashTable) forAll(f func(k, v uint32)) {
	numSlots := uint32(cap(iht.ba) / SlotSize)
	for slot := uint32(0); slot < numSlots; slot++ {
		byteLoc := slot * SlotSize
		if iht.ba[byteLoc] == FullSlot {
			k := readUInt32(iht.ba, byteLoc+1)
			v := readUInt32(iht.ba, byteLoc+5)
			f(k, v)
		}
	}
}

func (iht *IntHashTable) grow() {
	// double size of buffer and add all key values
	fmt.Println("Rehashing...")
	newBytes := make([]byte, 2*cap(iht.ba))
	inserter := func(k, v uint32) {
		put(newBytes, k, v)
	}
	iht.forAll(inserter)
	iht.ba = newBytes
}

func (iht *IntHashTable) Put(k, v uint32) {
	e := put(iht.ba, k, v)
	if e != nil {
		iht.grow()
		iht.Put(k, v)
	}
}

func (iht *IntHashTable) Get(k uint32) (uint32, bool) {
	numSlots := uint32(cap(iht.ba) / SlotSize)
	h := hash(k)
	slot := h % numSlots
	slotByteLoc := slot * SlotSize
	firstSlot := slot
	for {
		// slotByteLoc wraps around
		slotByteLoc = (slot % numSlots) * SlotSize
		if iht.ba[slotByteLoc] == FullSlot {
			readKey := readUInt32(iht.ba, slotByteLoc+1)
			if readKey == k {
				v := readUInt32(iht.ba, slotByteLoc+5)
				return v, true
			} else {
				// slot is monotonically increasing
				slot += 1
				if slot-firstSlot > numSlots/ScanFactor {
					panic(errors.New("Read scanned more than expected items, this should never happen."))
				}
			}
		} else {
			return 0, false
		}
	}
}

// note that this does not provide a copy, this is unsafe but exposed for
// performance reasons.
func (iht *IntHashTable) bytes() []byte {
	return iht.ba
}
