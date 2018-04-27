package milter

import (
	"bytes"
	"strings"
)

/* NULL terminator */
const NULL = "\x00"

/* DecodeCStrings splits c style strings into golang slice */
func DecodeCStrings(data []byte) (str []string) {
	if len(data) == 0 {
		return nil
	}
	str = strings.Split(string(data), NULL)
	return str[:len(str) - 1]
}

/* ReadCString reads and returs c style string from []byte */
func ReadCString(data []byte) string {
	pos := bytes.IndexByte(data, 0)
	if pos == -1 {
		return string(data)
	}
	return string(data[0:pos])
}
