package jacl

import (
	"fmt"
)

// ------------------------------------------------------------
// FIELDS

type Fields map[string]interface{}

// F is a convenience for building a field map. The data
// must be in pairs, where the first element of a pair is a string.
func F(pairs ...interface{}) Fields {
	f := make(map[string]interface{})
	key := ""
	for _, p := range pairs {
		if key != "" {
			f[key] = p
			key = ""
		} else {
			skey, ok := p.(string)
			if !ok {
				panic(fmt.Errorf("field element must be string"))
			}
			key = skey
		}
	}
	return f
}
