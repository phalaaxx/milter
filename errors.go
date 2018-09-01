package milter

import (
	"errors"
)

// pre-defined errors
var (
	eCloseSession = errors.New("Stop current milter processing")
	eMacroNoData  = errors.New("Macro definition with no data")
)
