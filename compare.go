package jacl

import (
	"fmt"
)

// compare() compares two interface values.
// All element of a must be in b, but not vice versa. It does not
// handle custom types, but assumes the values have been reduced
// to primitives.
func compare(a, b interface{}) bool {
	ans, err := compareBasicTypes(a, b)
	if err == nil {
		return ans
	}
	switch av := a.(type) {
	case map[string]interface{}:
		if bv, ok := b.(map[string]interface{}); ok {
			return compareStringInterfaceMap(av, bv)
		}
		return false
	case []interface{}:
		if bv, ok := b.([]interface{}); ok {
			return compareInterfaceSlice(av, bv)
		}
		return false
	}
	return a == b
}

// compareBasicTypes() compares basic types.
func compareBasicTypes(a, b interface{}) (bool, error) {
	if a == nil && b == nil {
		return true, nil
	}
	switch av := a.(type) {
	case string:
		if bv, ok := b.(string); ok {
			return av == bv, nil
		}
	case bool:
		if bv, ok := b.(bool); ok {
			return av == bv, nil
		}
	}
	return false, fmt.Errorf("Can't compare %T with %T", a, b)
}

// compareStringInterfaceMap() compares two maps of string to interface.
func compareStringInterfaceMap(a, b map[string]interface{}) bool {
	if a == nil && b == nil {
		return true
	} else if a != nil && b == nil {
		return false
	} else if a == nil && b != nil {
		return false
	}
	for ak, av := range a {
		bv, ok := b[ak]
		if !ok {
			return false
		}
		if !compare(av, bv) {
			return false
		}
	}
	return true
}

// compareInterfaceSlice() compares two maps of string to interface.
func compareInterfaceSlice(a, b []interface{}) bool {
	if len(a) != len(b) {
		return false
	} else if a == nil {
		return true
	}
	for i, ae := range a {
		if !compare(ae, b[i]) {
			return false
		}
	}
	return true
}
