package milter

import (
	"bufio"
	"encoding/binary"
	"io"
	"net"
	"net/textproto"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestMilterSession_ReadPacket(t *testing.T) {
	c1, c2 := net.Pipe()

	message := RespAccept.Response()
	// Send Message
	buffer := bufio.NewWriter(c1)
	length := uint32(len(message.Data) + 1)
	binary.Write(buffer, binary.BigEndian, length)
	buffer.Write(append([]byte{message.Code}, message.Data...))
	go buffer.Flush()

	sess := &milterSession{sock: c2}
	msg, err := sess.ReadPacket()
	if err != nil {
		t.Errorf("Unexpecte Error")
	}
	if msg.Code != message.Code {
		t.Errorf("Action don't match expected %s got %s",
			string(message.Code), string(msg.Code))
	}

	// Send smaller packet
	buffer = bufio.NewWriter(c1)
	binary.Write(buffer, binary.BigEndian, uint32(2))
	buffer.Write([]byte{message.Code})
	go func() {
		buffer.Flush()
		time.Sleep(10 * time.Millisecond)
		c1.Close()
	}()
	_, err = sess.ReadPacket()
	if err == nil {
		t.Errorf("Read Packet should fail here !!!")
	}
}

func TestMilterSession_WritePacket(t *testing.T) {
	c1, c2 := net.Pipe()

	message := RespAccept.Response()

	response := make(chan Message)
	go func() {
		var length uint32
		binary.Read(c1, binary.BigEndian, &length)
		data := make([]byte, length)
		io.ReadFull(c1, data)
		// prepare response data
		response <- Message{
			Code: data[0],
			Data: data[1:],
		}
	}()

	sess := &milterSession{sock: c2}
	err := sess.WritePacket(message)
	if err != nil {
		t.Errorf("Unexpected WritePacket %s", err)
	}

	msg := <-response
	if msg.Code != message.Code {
		t.Errorf("Action don't match expected %s got %s",
			string(message.Code), string(msg.Code))
	}

}

func TestMilterSession_ReadWrite(t *testing.T) {
	c1, c2 := net.Pipe()

	message := RespAccept.Response()

	sender := &milterSession{sock: c1}
	go func() {
		err := sender.WritePacket(message)
		if err != nil {
			t.Errorf("Unexpected WritePacket %s", err)
		}
	}()

	receiver := &milterSession{sock: c2}
	msg, err := receiver.ReadPacket()
	if err != nil {
		t.Errorf("Unexpected ReadPacket %s", err)
	}
	if msg.Code != message.Code {
		t.Errorf("Action don't match expected %s got %s",
			string(message.Code), string(msg.Code))
	}

}

type mocSession struct {
	Milter
}

func (moc *mocSession) Connect(host string, family string, port uint16, addr net.IP, m *Modifier) (Response, error) {
	return RespContinue, nil
}

func (moc *mocSession) Helo(name string, m *Modifier) (Response, error) {
	return RespContinue, nil
}

func (moc *mocSession) MailFrom(from string, m *Modifier) (Response, error) {
	return RespContinue, nil
}

func (moc *mocSession) RcptTo(rcptTo string, m *Modifier) (Response, error) {
	return RespContinue, nil
}

func (moc *mocSession) Header(name string, value string, m *Modifier) (Response, error) {
	return RespContinue, nil
}

func (moc *mocSession) Headers(h textproto.MIMEHeader, m *Modifier) (Response, error) {
	return RespContinue, nil
}

func (moc *mocSession) BodyChunk(chunk []byte, m *Modifier) (Response, error) {
	return RespContinue, nil
}
func (moc *mocSession) Body(m *Modifier) (Response, error) {
	return RespContinue, nil
}

func TestMilterSession_Process(t *testing.T) {

	sess := &milterSession{
		milter: &mocSession{},
	}
	tests := []struct {
		Message    *Message
		ShouldFail bool
	}{
		{
			Message:    SimpleResponse('A').Response(),
			ShouldFail: false,
		},
		{
			Message:    SimpleResponse('B').Response(),
			ShouldFail: false,
		},
		{
			Message:    NewResponse('C', []byte("server.domain.com\x004\x00\x01\x01127.0.0.1\x00")).Response(),
			ShouldFail: false,
		},
		{
			Message:    NewResponse('C', []byte("server.domain.com\x004\x00")).Response(),
			ShouldFail: false,
		},
		{
			Message:    NewResponse('C', []byte("server.domain.com\x006\x00\x01\x01IPv6:2001:db8:1234:ffff:ffff:ffff:ffff:ffff")).Response(),
			ShouldFail: false,
		},
		{
			Message:    NewResponse('D', []byte("key\x00value\x00")).Response(),
			ShouldFail: false,
		},
		{
			Message:    SimpleResponse('E').Response(),
			ShouldFail: false,
		},
		{
			Message:    NewResponse('H', []byte("mail.domain.com\x00")).Response(),
			ShouldFail: false,
		},
		{
			Message:    NewResponse('L', []byte("From\x00user@domain.com\x00")).Response(),
			ShouldFail: false,
		},
		{
			Message:    NewResponse('M', []byte("user@domain.com\x00")).Response(),
			ShouldFail: false,
		},
		{
			Message:    SimpleResponse('N').Response(),
			ShouldFail: false,
		},
		{
			Message:    SimpleResponse('O').Response(),
			ShouldFail: false,
		},
		{
			Message:    SimpleResponse('Q').Response(),
			ShouldFail: true,
		},
		{
			Message:    SimpleResponse('R').Response(),
			ShouldFail: false,
		},
		{
			Message:    SimpleResponse('T').Response(),
			ShouldFail: false,
		},
		{
			Message:    SimpleResponse('Z').Response(),
			ShouldFail: true,
		},
	}

	for _, test := range tests {
		_, err := sess.Process(test.Message)
		if test.ShouldFail {
			assert.Error(t, err)
		} else {
			assert.NoError(t, err)
		}
	}

}

func TestMilterSession_HandleMilterCommands(t *testing.T) {

}
