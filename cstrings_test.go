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
func TestDecodeCStringsTwoNull(t *testing.T) {
	strs := decodeCStrings([]byte{})
	if len(strs) > 0 {
		t.Errorf("Expecting empty slice")
	}

	strs = decodeCStrings([]byte("Gopher\x00Academy\x00Test\x00\x00"))
	if len(strs) != 4 {
		t.Errorf("Expecting slice with 2 strings got %#v", strs)
	}
	if strs[0] != "Gopher" {
		t.Errorf("Expecting `Gopher` got %s", strs[0])
	}
	if strs[1] != "Academy" {
		t.Errorf("Expecting `Academy` got %s", strs[1])
	}
	if strs[2] != "Test" {
		t.Errorf("Expecting `Test` got %s", strs[2])
	}
	if strs[3] != "" {
		t.Errorf("Expecting `` got %s", strs[3])
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
