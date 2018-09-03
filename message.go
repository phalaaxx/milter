package milter

// Message represents a command sent from milter client
type Message struct {
	Code byte
	Data []byte
}

// Define milter response codes
const (
	accept    = 'a'
	continue_ = 'c'
	discard   = 'd'
	reject    = 'r'
	tempFail  = 't'
)
