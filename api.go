package jacl

// A Cmper compares the values of two types.
type Cmper interface {
	// Answer true if b contains all the values of a.
	// Note this is not symmetric: b can contain more
	// values then a and it is still true.
	Cmp(b interface{}) (bool, error)

	// Answer a unique key so I can be reinstantiated after marshalling.
	FactoryKey() string
}

// Cmp() constructs a new comparison object. It can be used
// against a single item. The item must resolve to a map
// of string -> interface{}.
func Cmp(a interface{}) Cmper {
	return singleCmp{A: a}
}

// Cmps() constructs a new comparison object. It can be used
// against a slice of items. Each item in the slice must resolve
// to a map of string -> interface{}.
// key is an optional key to identify equivalent items between the
// slices being compared. If it is empty, then the slices are assumed
// to be in sorted order. If it's not empty, then it is used as a lookup
// to find items to compare in each slice. It can be either a string,
// if a single value is used to uniquely identify items, or a slice
// of strings, if multiple values are required. Panic if the key is invalid.
func Cmps(_key interface{}, _a ...interface{}) Cmper {
	key := convertKey(_key)
	var a []interface{}
	var fn []FuncFactory
	for _, ai := range _a {
		switch ait := ai.(type) {
		case CmpsFunc:
			fn = append(fn, FuncFactory{Fn: ait})
		default:
			a = append(a, ai)
		}
	}
	return sliceCmp{Keys: key, A: a, Fn: fn}
}

// SizeIs() can be passed as one of the values to Cmps. It is a special
// comparison function: Error if the result size does not match.
func SizeIs(size int) interface{} {
	return &sizeis{Size: size}
}

// Compare() compares two interface values. Clients can use it
// directly instead of going through the Cmper interface.
func Compare(a, b interface{}) bool {
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
