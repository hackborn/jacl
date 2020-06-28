package jacl

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
	return false, errIncomparable
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
		if !Compare(av, bv) {
			return false
		}
	}
	return true
}
