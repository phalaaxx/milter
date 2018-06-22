package milter

import (
	"errors"
)

// pre-defined errors
var (
	ECloseSession = errors.New("Stop current milter processing")
	EMacroNoData  = errors.New("Macro definition with no data")
)
