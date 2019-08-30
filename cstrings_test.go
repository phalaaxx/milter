package milter

import (
	"testing"
)

func TestDecodeCStrings(t *testing.T) {
	strs := decodeCStrings([]byte{})
	if len(strs) > 0 {
		t.Errorf("Expecting empty slice")
	}

	strs = decodeCStrings([]byte("Gopher\x00Academy\x00"))
	if len(strs) != 2 {
		t.Errorf("Expecting slice with 2 strings got %#v", strs)
	}
	if strs[0] != "Gopher" {
		t.Errorf("Expecting `Gopher` got %s", strs[0])
	}
	if strs[1] != "Academy" {
		t.Errorf("Expecting `Academy` got %s", strs[1])
	}
}

func TestReadCString(t *testing.T) {
	str := readCString([]byte(""))
	if len(str) > 0 {
		t.Errorf("Expecting 0 byte")
	}

	str = readCString([]byte("GO\x00"))
	if len(str) == 0 {
		t.Errorf("Expecting len 2")
	}
	if str != "GO" {
		t.Errorf("Expecting `GO` got %s", str)
	}
}
