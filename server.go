/* RunServer provides a convinient way to start milter server main loop */
package milter

import (
	"net"
)

/* Run server listener */
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
