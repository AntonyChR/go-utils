// bytes package provides functions for byte manipulation.
package bytes

func uint16ToBytes(n uint16)[]byte{
	p1 := uint8(n >> 8)
	p2 := uint8(n & 0xFF)
	return []byte{p1,p2}
}

