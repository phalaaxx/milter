package milter

import (
	"bufio"
	"encoding/binary"
	"io"
	"net"
	"testing"
	"time"
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

func TestMilterSession_Process(t *testing.T) {
	sess := &milterSession{}

	_, err := sess.Process(SimpleResponse('A').Response())
	if err != nil {
		t.Errorf("Unexpected Process %s", err)
	}
	/*
		_,err = sess.Process(SimpleResponse('B').Response())
		if err != nil {
			t.Errorf("Unexpected Process %s", err)
		}
	*/
	/*
		_,err = sess.Process(SimpleResponse('C').Response())
		if err != nil {
			t.Errorf("Unexpected Process %s", err)
		}
	*/

}

func TestMilterSession_HandleMilterCommands(t *testing.T) {

}
