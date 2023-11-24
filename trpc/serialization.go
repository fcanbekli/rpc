package main

import "encoding/binary"

func SerializeInt(data []byte, value int, left int, right int) {
	// Use encoding/binary to serialize the integer into the byte slice
	binary.BigEndian.PutUint64(data[left:right], uint64(value))
}

func DeserializeInt(data []byte, start int, end int) int {
	// Use encoding/binary to deserialize the integer from the byte slice
	return int(binary.BigEndian.Uint64(data[start:end]))
}
