package milter

// Message represents a command sent from milter client
type Message struct {
	Code byte
	Data []byte
}

// Define milter response codes
const (
	Accept   = 'a'
	Continue = 'c'
	Discard  = 'd'
	Reject   = 'r'
	TempFail = 't'
)
