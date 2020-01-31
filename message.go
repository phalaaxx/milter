package milter

// Message represents a command sent from milter client
type Message struct {
	Code byte
	Data []byte
}

// Define milter response codes
const (
	rspAccept   = 'a'
	rspContinue = 'c'
	rspDiscard  = 'd'
	rspReject   = 'r'
	rspTempFail = 't'
)
