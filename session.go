package milter

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"io"
	"log"
	"net"
	"net/textproto"
	"strings"
)

const (
	// negotiation actions
	OptAddHeader    = 0x01
	OptChangeBody   = 0x02
	OptAddRcpt      = 0x04
	OptRemoveRcpt   = 0x08
	OptChangeHeader = 0x10
	OptQuarantine   = 0x20

	// undesired protocol content
	OptNoConnect  = 0x01
	OptNoHelo     = 0x02
	OptNoMailFrom = 0x04
	OptNoRcptTo   = 0x08
	OptNoBody     = 0x10
	OptNoHeaders  = 0x20
	OptNoEOH      = 0x40
)

/* MilterSession is keeps session state during MTA communication */
type MilterSession struct {
	Actions  uint32
	Protocol uint32
	Sock     io.ReadWriteCloser
	Headers  textproto.MIMEHeader
	Macros   map[string]string
	Milter   Milter
}

/* ReadPacket reads incoming milter packet */
func (c *MilterSession) ReadPacket() (*Message, error) {
	// read packet length
	var length uint32
	if err := binary.Read(c.Sock, binary.BigEndian, &length); err != nil {
		return nil, err
	}

	// read packet data
	data := make([]byte, length)
	if _, err := io.ReadFull(c.Sock, data); err != nil {
		return nil, err
	}

	// prepare response data
	message := Message{
		Code: data[0],
		Data: data[1:],
	}

	return &message, nil
}

/* WritePacket sends a milter response packet to socket stream */
func (m *MilterSession) WritePacket(msg *Message) error {
	buffer := bufio.NewWriter(m.Sock)

	// calculate and write response length
	length := uint32(len(msg.Data) + 1)
	if err := binary.Write(buffer, binary.BigEndian, length); err != nil {
		return err
	}

	// write response code
	if err := buffer.WriteByte(msg.Code); err != nil {
		return err
	}

	// write response data
	if _, err := buffer.Write(msg.Data); err != nil {
		return err
	}

	// flush data to network socket stream
	if err := buffer.Flush(); err != nil {
		return err
	}

	return nil
}

/* Process incoming milter command */
func (m *MilterSession) Process(msg *Message) (Response, error) {
	switch msg.Code {
	case 'A':
		// abort current message and start over
		m.Headers = nil
		m.Macros = nil
		// do not send response
		return nil, nil

	case 'B':
		// body chunk
		return m.Milter.BodyChunk(msg.Data, NewModifier(m))

	case 'C':
		// new connection, get hostname
		Hostname := ReadCString(msg.Data)
		msg.Data = msg.Data[len(Hostname)+1:]
		// get protocol family
		ProtocolFamily := msg.Data[0]
		msg.Data = msg.Data[1:]
		// get port
		var Port uint16
		if ProtocolFamily == '4' || ProtocolFamily == '6' {
			if len(msg.Data) < 2 {
				return RespTempFail, nil
			}
			Port = binary.BigEndian.Uint16(msg.Data)
			msg.Data = msg.Data[2:]
		}
		// get address
		Address := ReadCString(msg.Data)
		// convert address and port to human readable string
		family := map[byte]string{
			'U': "unknown",
			'L': "unix",
			'4': "tcp4",
			'6': "tcp6",
		}
		// run handler and return
		return m.Milter.Connect(
			Hostname,
			family[ProtocolFamily],
			Port,
			net.ParseIP(Address),
			NewModifier(m))

	case 'D':
		// define macros
		m.Macros = make(map[string]string)
		// convert data to golang strings
		data := DecodeCStrings(msg.Data[1:])
		if len(data) != 0 {
			// store data in a map
			for i := 0; i < len(data); i += 2 {
				m.Macros[data[i]] = data[i+1]
			}
		}
		// do not send response
		return nil, nil

	case 'E':
		// call and return milter handler
		return m.Milter.Body(NewModifier(m))

	case 'H':
		// helo command
		name := strings.TrimSuffix(string(msg.Data), NULL)
		return m.Milter.Helo(name, NewModifier(m))

	case 'L':
		// make sure Headers is initialized
		if m.Headers == nil {
			m.Headers = make(textproto.MIMEHeader)
		}
		// add new header to headers map
		HeaderData := DecodeCStrings(msg.Data)
		if len(HeaderData) == 2 {
			m.Headers.Add(HeaderData[0], HeaderData[1])
			// call and return milter handler
			return m.Milter.Header(HeaderData[0], HeaderData[1], NewModifier(m))
		}

	case 'M':
		// envelope from address
		envfrom := ReadCString(msg.Data)
		return m.Milter.MailFrom(strings.Trim(envfrom, "<>"), NewModifier(m))

	case 'N':
		// end of headers
		return m.Milter.Headers(m.Headers, NewModifier(m))

	case 'O':
		// ignore request and prepare response buffer
		buffer := new(bytes.Buffer)
		// prepare response data
		for _, value := range []uint32{2, m.Actions, m.Protocol} {
			if err := binary.Write(buffer, binary.BigEndian, value); err != nil {
				return nil, err
			}
		}
		// build and send packet
		return NewResponse('O', buffer.Bytes()), nil

	case 'Q':
		// client requested session close
		return nil, ECloseSession

	case 'R':
		// envelope to address
		envto := ReadCString(msg.Data)
		return m.Milter.RcptTo(strings.Trim(envto, "<>"), NewModifier(m))

	case 'T':
		// data, ignore

	default:
		// print error and close session
		log.Printf("Unrecognized command code: %c", msg.Code)
		return nil, ECloseSession
	}

	// by default continue with next milter message
	return RespContinue, nil
}

/* process all milter commands in the same connection */
func (m *MilterSession) HandleMilterCommands() {
	// close session socket on exit
	defer m.Sock.Close()

	for {
		// ReadPacket
		msg, err := m.ReadPacket()
		if err != nil {
			if err != io.EOF {
				log.Printf("Error reading milter command: %v", err)
			}
			return
		}

		// process command
		resp, err := m.Process(msg)
		if err != nil {
			if err != ECloseSession {
				// log error condition
				log.Printf("Error performing milter command: %v", err)
			}
			return
		}

		// ignore empty responses
		if resp != nil {
			// send back response message
			if err = m.WritePacket(resp.Response()); err != nil {
				log.Printf("Error writing packet: %v", err)
				return
			}

			if !resp.Continue() {
				return
			}

		}
	}
}
