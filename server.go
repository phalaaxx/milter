// A Go library for milter support
package milter

import (
	"errors"
	"fmt"
	"net"
	"sync"
)

var defaultServer Server

// MilterInit initializes milter options
// multiple options can be set using a bitmask
type MilterInit func() (Milter, OptAction, OptProtocol)

// RunServer provides a convenient way to start a milter server
// Handlers provide way to handle errors from panics
// With nil handlers panics not recovered
func RunServer(server net.Listener, init MilterInit, handlers ...func(error)) error {
	defaultServer.Listener = server
	defaultServer.MilterFactory = init
	defaultServer.ErrHandlers = handlers
	return defaultServer.RunServer()
}

// Close server listener and wait worked process
func Close() (err error) {
	return defaultServer.Close()
}

// Milter server for handling and processing incoming connections
// support panic handling via ErrHandler
// couple of func(error) could be provided for handling error
type Server struct {
	Listener      net.Listener
	MilterFactory MilterInit
	ErrHandlers   []func(error)
	sync.WaitGroup
}

// For graceful shutdown
// Stop accepting new connections
// And wait until processing connections ends
func (s *Server) Close() (err error) {
	if s.Listener != nil {
		err = s.Listener.Close()
	}
	s.Wait()
	return err
}

// Start milter server via provided listener
func (s *Server) RunServer() error {
	if s.Listener == nil {
		return errors.New("no listen addr specified")
	}

	for {
		// accept connection from client
		conn, err := s.Listener.Accept()
		if err != nil {
			return err
		}

		go func() {
			defer handlePanic(s.ErrHandlers)
			defer s.Done()
			s.Add(1)
			s.handleCon(conn)
		}()
	}
}

// Handle incoming connections
func (s *Server) handleCon(conn net.Conn) {
	// create milter object
	milter, actions, protocol := s.MilterFactory()
	session := milterSession{
		actions:  actions,
		protocol: protocol,
		sock:     conn,
		milter:   milter,
	}
	// handle connection commands
	session.HandleMilterCommands()
}

// Recover panic from session and call handle with occurred error
// If no any handle provided panics will not recovered
func handlePanic(handlers []func(error)) {
	var err error

	if handlers == nil {
		return
	}

	r := recover()
	switch r.(type) {
	case nil:
		return
	case error:
		err = r.(error)
	default:
		err = errors.New(fmt.Sprint(r))
	}
	for _, f := range handlers {
		f(err)
	}
}
