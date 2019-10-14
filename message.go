package milter

// Message represents a command sent from milter client
type Message struct {
	Code byte
	Data []byte
}

// Define milter response codes
const (
	acceptAction   = 'a'
	continueAction = 'c'
	discardAction  = 'd'
	rejectAction   = 'r'
	tempFailAction = 't'
)
