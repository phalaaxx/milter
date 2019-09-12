package milter

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestModifier_AddRecipient(t *testing.T) {
	recipient := "history@domain.com"

	m := &Modifier{
		writePacket: func(m *Message) error {
			assert.Equal(t, m.Code, byte('+'))
			assert.Equal(t, readCString(m.Data), "<"+recipient+">")
			return nil
		},
	}

	err := m.AddRecipient(recipient)
	assert.NoError(t, err)
}

func TestModifier_DeleteRecipient(t *testing.T) {
	recipient := "history@domain.com"

	m := &Modifier{
		writePacket: func(m *Message) error {
			assert.Equal(t, m.Code, byte('-'))
			assert.Equal(t, readCString(m.Data), "<"+recipient+">")
			return nil
		},
	}

	err := m.DeleteRecipient(recipient)
	assert.NoError(t, err)
}

func TestModifier_ReplaceBody(t *testing.T) {
	body := []byte("Replaced Message body")
	m := &Modifier{
		writePacket: func(m *Message) error {
			assert.Equal(t, m.Code, byte('b'))
			assert.Equal(t, m.Data, body)
			return nil
		},
	}

	err := m.ReplaceBody(body)
	assert.NoError(t, err)
}

func TestModifier_AddHeader(t *testing.T) {
	m := &Modifier{
		writePacket: func(m *Message) error {
			return nil
		},
	}

	err := m.AddHeader("X-Mail-Filter", "milter-daemon")
	assert.NoError(t, err)
}

func TestModifier_Quarantine(t *testing.T) {
	reason := "history@domain.com"

	m := &Modifier{
		writePacket: func(m *Message) error {
			assert.Equal(t, readCString(m.Data), reason)
			return nil
		},
	}

	err := m.Quarantine(reason)
	assert.NoError(t, err)
}

func TestModifier_ChangeHeader(t *testing.T) {

	name, value := "X-Mail-Filter", "milter-daemon"

	m := &Modifier{
		writePacket: func(m *Message) error {
			assert.Equal(t, m.Code, byte('m'))
			//assert.Equal(t,m.Data,body)
			return nil
		},
	}

	err := m.ChangeHeader(1, name, value)
	assert.NoError(t, err)
}
