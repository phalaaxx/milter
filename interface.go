package milter

import (
	"net"
	"net/textproto"
)

/* Milter is an interface for milter callback handlers */
type Milter interface {
	// Connect is called to provide SMTP connection data for incoming message
	//   supress with NoConnect
	Connect(string, string, uint16, net.IP, *Modifier) (Response, error)

	// Helo is called to process any HELO/EHLO related filters
	//   supress with NoHelo
	Helo(string, *Modifier) (Response, error)

	// MailFrom is called to process filters on envelope FROM address
	//   supress with NoMailForm
	MailFrom(string, *Modifier) (Response, error)

	// RcptTo is called to process filters on envelope TO address
	//   supress with NoRcptTo
	RcptTo(string, *Modifier) (Response, error)

	// Header is called once for each header in incoming message
	//   supress with NoHeaders
	Header(string, string, *Modifier) (Response, error)

	// Headers is called when all message headers have been processed
	//   supress with NoHeaders
	Headers(textproto.MIMEHeader, *Modifier) (Response, error)

	// BodyChunk is called to process next message body chunk data (up to 64KB in size)
	//   supress with NoBody
	BodyChunk([]byte, *Modifier) (Response, error)

	// Body is called when entire message body has been passed to milter
	//   supress with NoBody
	Body(*Modifier) (Response, error)
}
