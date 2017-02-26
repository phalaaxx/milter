package milter

import (
	"net"
)

/* MilterInit is a function that initializes milter object options */
type MilterInit func() (Milter, uint32, uint32)

/* RunServer provides a convinient way to start milter server */
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
