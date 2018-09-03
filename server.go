// A Go library for milter support
package milter

import (
	"net"
)

// MilterInit initializes milter options
// multiple options can be set using a bitmask
type MilterInit func() (Milter, OptAction, OptProtocol)

// RunServer provides a convenient way to start a milter server
func RunServer(server net.Listener, init MilterInit) error {
	for {
		// accept connection from client
		client, err := server.Accept()
		if err != nil {
			return err
		}
		// create milter object
		milter, actions, protocol := init()
		session := MilterSession{
			Actions:  actions,
			Protocol: protocol,
			Sock:     client,
			Milter:   milter,
		}
		// handle connection commands
		go session.HandleMilterCommands()
	}
}
