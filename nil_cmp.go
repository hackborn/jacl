package jacl

import (
	"fmt"
	"reflect"
)

// ------------------------------------------------------------
// NIL-CMP

// nilCmp compares a single item to another.
type nilCmp struct {
}

func (c nilCmp) Cmp(b interface{}) error {
	if !isNilInterface(b) {
		return newComparisonError(fmt.Sprintf(haveWantFmt, toJson(b), `nil`))
	}
	return nil
}

func (c nilCmp) SerializeKey() string {
	return nilCmpFactoryKey
}

func isNilInterface(i interface{}) bool {
	if i == nil {
		return true
	}
	switch reflect.TypeOf(i).Kind() {
	case reflect.Ptr, reflect.Map, reflect.Array, reflect.Chan, reflect.Slice:
		return reflect.ValueOf(i).IsNil()
	}
	return false
}
