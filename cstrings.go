package milter

import (
	"bytes"
	"strings"
)

// NULL terminator
const null = "\x00"

// DecodeCStrings splits a C style strings into a Go slice
func decodeCStrings(data []byte) []string {
	if len(data) == 0 {
		return nil
	}
	str := strings.Split(string(data), null)
	return str[:len(str)-1]
}

// ReadCString reads and returns a C style string from []byte
func readCString(data []byte) string {
	pos := bytes.IndexByte(data, 0)
	if pos == -1 {
		return string(data)
	}
	return string(data[0:pos])
}
