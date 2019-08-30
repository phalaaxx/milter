package milter

import (
	"testing"
)

func TestResponse(t *testing.T) {

}

func TestNewResponse(t *testing.T) {
	if NewResponse(acceptAction, []byte{}).Continue() != false {
		t.Fail()
	}

	rData := []byte{0x01, 0x02, 0x03}
	response := NewResponse(continueAction, rData)
	if response.Continue() == false {
		t.Fail()
	}

	if response.Response().Code != continueAction {
		t.Errorf("Response action don't match expected %s got %s",
			string(continueAction), string(response.Response().Code))
	}

	if len(response.Response().Data) != len(rData) {
		t.Errorf("Response data length don't match expected %d , got %d",
			len(rData), len(response.Response().Data))
	}
}
