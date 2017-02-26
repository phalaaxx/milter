package milter

import (
	"net"
)

/* RunServer provides a convinient way to start milter server */
func RunServer(server net.Listener, MilterInit func() (Milter, uint32, uint32)) error {
	for {
		// accept connection from client
		client, err := server.Accept()
		if err != nil {
			return err
		}
		// create milter object
		milter, actions, protocol := MilterInit()
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
