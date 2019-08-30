package milter

import (
	"errors"
	"net"
	"testing"
	"time"
)

type mocMilter struct {
	Milter
}

type mocListener struct {
	Times    int
	Listener net.Listener
}

func (m mocListener) Accept() (net.Conn, error) {
	if m.Listener == nil {
		return nil, errors.New("Failed")
	}
	return m.Listener.Accept()
}

func (m mocListener) Close() error {
	return nil
}

func (m mocListener) Addr() net.Addr {
	return nil
}

func TestRunServer(t *testing.T) {

	milterInit := func() (Milter, OptAction, OptProtocol) {
		return &mocMilter{},
			OptAddHeader,
			OptNoConnect | OptNoHelo | OptNoMailFrom | OptNoRcptTo
	}

	ln, err := net.Listen("tcp4", "127.0.0.1:")
	if err != nil {
		t.Errorf("Listen")
	}

	listener := &mocListener{Listener: ln}

	go func(address string) {
		time.Sleep(10 * time.Millisecond)
		cli, err := net.Dial("tcp4", address)
		if err != nil {
			t.Errorf("Dial %s", err)
		}
		sess := milterSession{
			sock: cli,
		}

		// Send invalid Command
		sess.WritePacket(SimpleResponse('Q').Response())
		// End session
		cli.Close()

		// New Session
		listener.Listener = nil

		cli, err = net.Dial("tcp4", address)
		if err != nil {
			t.Errorf("Dial %s", err)
		}
		sess.sock = cli
		//sess.WritePacket(SimpleResponse('X').Response())
		cli.Close()
	}(ln.Addr().String())

	RunServer(listener, milterInit)
}
