package jacl

import (
	"fmt"
)

// ------------------------------------------------------------
// TESTING

func equalErr(a, b error) bool {
	if a == b {
		return true
	} else if a == nil {
		return false
	} else if b == nil {
		return false
	} else {
		return a.Error() == b.Error()
	}
}

// ------------------------------------------------------------
// CONST and VAR

var (
	errIncomparable = fmt.Errorf("Incomparable")
)
