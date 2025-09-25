// bytes package provides functions for byte manipulation.
package bytes

func uint16ToBytes(n uint16) []byte {
	return []byte{byte(n >> 8), byte(n)}
}
