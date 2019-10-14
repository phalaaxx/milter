package milter

import (
	"testing"
)

func TestSimpleResponse(t *testing.T) {
	response := SimpleResponse(continueAction)
	if response.Continue() == false {
		t.Fail()
	}

	if response.Response().Code != continueAction {
		t.Errorf("Response action don't match expected %s got %s",
			string(continueAction), string(response.Response().Code))
	}

	if len(response.Response().Data) != 0 {
		t.Errorf("Response data length don't match expected %d , got %d",
			0, len(response.Response().Data))
	}
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

func TestNewResponseStr(t *testing.T) {
	rData := "Message"
	response := NewResponseStr(continueAction, rData)

	if response.Response().Code != continueAction {
		t.Errorf("Response action don't match expected %s got %s",
			string(continueAction), string(response.Response().Code))
	}

	if len(response.Response().Data) != len(rData)+1 {
		t.Errorf("Response data length don't match expected %d , got %d",
			len(rData)+1, len(response.Response().Data))
	}
}
