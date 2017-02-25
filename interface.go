package milter

import (
	"net/textproto"
)

/* Milter is an interface for milter callback handlers */
type Milter interface {
	Connect(string, string, string, *Modifier) (Response, error)
	Helo(string, *Modifier) (Response, error)
	MailFrom(string, *Modifier) (Response, error)
	RcptTo(string, *Modifier) (Response, error)
	Header(string, string, *Modifier) (Response, error)
	Headers(textproto.MIMEHeader, *Modifier) (Response, error)
	BodyChunk([]byte, *Modifier) (Response, error)
	Body([]byte, *Modifier) (Response, error)
}
