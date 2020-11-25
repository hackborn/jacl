package jacl

// A Cmper compares the values of two types.
type Cmper interface {
	// Answer nil if b contains all the values of a,
	// an error otherwise. Note this is not symmetric:
	// b can contain more values than a.
	// The API is designed for testing scenarios, where
	// you generally want to know pass/fail: If you
	// need more information (for example, whether the
	// test ran fine but the comparison failed, vs as
	// an error that occurred attempting the comparison),
	// examine the error type, which will be either a
	// ComparisonError or EvaluationError.
	Cmp(b interface{}) error
}

// Cmp() constructs a new comparison object. It can be used
// against a single item. The item must resolve to a map
// of string -> interface{}.
func Cmp(a interface{}) Cmper {
	return singleCmp{A: a}
}

// Cmps() constructs a new comparison object to be used against a
// slice of items. Each item in the slice must resolve
// to a map of string -> interface{}.
//
// Additional functionality is available via cmps funcs. See below.
func Cmps(_a ...interface{}) Cmper {
	var key []string
	var a []interface{}
	var fn []FuncFactory
	for _, ai := range _a {
		switch ait := ai.(type) {
		case keyFn:
			key = ait.Keys
		case *keyFn:
			key = ait.Keys
		case CmpsFunc:
			fn = append(fn, FuncFactory{Fn: ait})
		default:
			a = append(a, ai)
		}
	}
	return sliceCmp{Keys: key, A: a, Fn: fn}
}

// NilCmp() constructs a new comparison object that fails
// if the comparison is not nil.
// of string -> interface{}.
func CmpNil() Cmper {
	return nilCmp{}
}

// ------------------------------------------------------------
// CMPS FUNCS

// Key() can be passed as one of the values to Cmps(). It is a special
// matching function: It defines what keys are used to determine identity
// between the two slices being compared. This can be used to compare
// slices of unequal size, or slices in different orders.
func Key(v ...string) interface{} {
	return &keyFn{Keys: v}
}

// SizeIs() can be passed as one of the values to Cmps(). It is a special
// comparison function: Error if the result size does not match.
func SizeIs(size int) interface{} {
	return &sizeisFn{Size: size}
}
