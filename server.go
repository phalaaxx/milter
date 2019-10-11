// Package milter A Go library for milter support
package milter

import (
	"errors"
	"net"
	"sync"
	"time"
)

// An uninteresting service.
type Server struct {
	milter   Milter
	actions  OptAction
	protocol OptProtocol
	doneChan chan struct{}
	mu       sync.Mutex
}

func New(milter Milter, actions OptAction, protocol OptProtocol) *Server {
	return &Server{
		actions:  actions,
		protocol: protocol,
		milter:   milter,
	}
}

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
			actions:  srv.actions,
			protocol: srv.protocol,
			sock:     client,
			milter:   srv.milter,
		}

		wg.Add(1)
		go func() {
			wg.Done()
			session.HandleMilterCommands()
		}()
	}
}

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
