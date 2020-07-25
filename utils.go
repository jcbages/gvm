package main

import (
	"encoding/binary"
	"io"
)

// U1 represents an unsigned 1 byte int
type U1 uint8

// U2 represents an unsigned 2 byte int
type U2 uint16

// U4 represents an unsigned 4 byte int
type U4 uint32

// ReadU1 reads a u1 int
func ReadU1(r io.Reader) U1 {
	result := readBytes(r, 1)
	return U1(result[0])
}

// ReadU2 reads a u2 int
func ReadU2(r io.Reader) U2 {
	result := readBytes(r, 2)
	return U2(binary.BigEndian.Uint16(result))
}

// ReadU4 reads a u4 int
func ReadU4(r io.Reader) U4 {
	result := readBytes(r, 4)
	return U4(binary.BigEndian.Uint32(result))
}

func readBytes(r io.Reader, size int) []byte {
	result := make([]byte, size)
	received, err := r.Read(result)
	if err != nil || received != size {
		panic(err)
	}
	return result
}
