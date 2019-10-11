// Package milter A Go library for milter support
package milter

import (
	"errors"
	"net"
	"sync"
	"time"
)

// Server Milter Server
type Server struct {
	Milter   Milter
	Actions  OptAction
	Protocol OptProtocol
	doneChan chan struct{}
	mu       sync.Mutex
}

// New return Milter Server instance
func New(milter Milter, actions OptAction, protocol OptProtocol) *Server {
	return &Server{
		Actions:  actions,
		Protocol: protocol,
		Milter:   milter,
	}
}

// Serve accepts incoming connections on the Listener l, creating a
// new service goroutine for each. The service goroutines read requests and
// then call milterSession.Handler to reply to them.
func (srv *Server) Serve(listener net.Listener) error {
	wg := sync.WaitGroup{}
	for {

		select {
		case <-srv.getDoneChan():
			listener.Close()
			wg.Wait()
			return nil
		default:
		}
		switch lstnr := listener.(type) {
		case *net.TCPListener:
			lstnr.SetDeadline(time.Now().Add(1e9))
		case *net.UnixListener:
			lstnr.SetDeadline(time.Now().Add(1e9))
		default:
			errors.New("Unknown listener type")
		}

		// accept connection from client
		client, err := listener.Accept()
		if nil != err {
			if opErr, ok := err.(*net.OpError); ok && opErr.Timeout() {
				continue
			}
			return err
		}

		// create milter object
		session := milterSession{
			actions:  srv.Actions,
			protocol: srv.Protocol,
			sock:     client,
			milter:   srv.Milter,
		}

		wg.Add(1)
		go func() {
			wg.Done()
			session.Handle()
		}()
	}
}

// Shutdown gracefully shuts down the server without interrupting any
// active connections. Shutdown works by first closing all open
// listeners, then closing all idle connections, and then waiting
// indefinitely for connections to return to idle and then shut down.
func (srv *Server) Shutdown() {
	srv.mu.Lock()
	defer srv.mu.Lock()

	ch := srv.getDoneChanLocked()
	select {
	case <-ch:
		// Already closed. Don't close again.
	default:
		// Safe to close here. We're the only closer, guarded
		// by srv.mu.
		close(ch)
	}

}

func (srv *Server) getDoneChan() <-chan struct{} {
	srv.mu.Lock()
	defer srv.mu.Unlock()
	return srv.getDoneChanLocked()
}

func (srv *Server) getDoneChanLocked() chan struct{} {
	if srv.doneChan == nil {
		srv.doneChan = make(chan struct{})
	}
	return srv.doneChan
}
