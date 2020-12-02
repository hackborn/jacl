package jacl

import (
	"fmt"
)

// ------------------------------------------------------------
// SLICE-CMP

// sliceCmp compares a list of items to another list.
// It has an optional key: If present, the key is used
// as an identifier to match items in each list. If absent,
// it assumes the list are in the correct order.
type sliceCmp struct {
	Keys []string      `json:"key,omitempty"`
	A    []interface{} `json:"a,omitempty"`
	Fn   []FuncFactory `json:"fn,omitempty"`
}

func (c sliceCmp) Cmp(b interface{}) error {
	if c.A == nil && b == nil {
		return nil
	}
	bslice := make([]interface{}, 0)
	err := toFromJson(b, &bslice)
	if err != nil {
		return newEvaluationError(err)
	}

	for _, fn := range c.Fn {
		err = fn.Eval(bslice)
		if err != nil {
			return newEvaluationError(err)
		}
	}
	// If we have functions but no data, we don't
	// perform a comparison. This handles the case where
	// the comparison is nothing but not-exists checks.
	if len(c.Fn) > 0 && len(c.A) < 1 {
		return nil
	}

	asrc, bsrc, err := c.convertToStringMaps(bslice)
	if err == nil {
		return c.cmpStringMaps(asrc, bsrc)
	}

	// If I couldn't convert to string maps, assume the slices
	// contain literals.
	return c.cmpSlices(c.A, bslice)
}

func (c sliceCmp) SerializeKey() string {
	return sliceCmpFactoryKey
}

func (c sliceCmp) addFn(_fn interface{}) sliceCmp {
	if fn, ok := _fn.(CmpsFunc); ok {
		c.Fn = append(c.Fn, FuncFactory{Fn: fn})
		return c
	}
	panic("unknown func")
}

func (c sliceCmp) cmpStringMaps(asrc, bsrc []map[string]interface{}) error {
	for i, av := range asrc {
		bv := c.find(c.Keys, i, av, bsrc)
		if bv == nil {
			return newComparisonError(fmt.Sprintf(haveWantFmt, toJson(bsrc), toJson(asrc)))
		}
		if !compare(av, bv) {
			return newComparisonError(fmt.Sprintf(haveWantFmt, toJson(bsrc), toJson(asrc)))
		}
	}
	return nil
}

func (c sliceCmp) cmpSlices(aslice, bslice []interface{}) error {
	if len(aslice) != len(bslice) {
		return newComparisonError(fmt.Sprintf(haveWantLengthFmt, len(bslice), len(aslice)))
	}
	for i, av := range aslice {
		if !compare(av, bslice[i]) {
			return newComparisonError(fmt.Sprintf(haveWantFmt, toJson(bslice), toJson(aslice)))
		}
	}
	return nil
}

func (c sliceCmp) find(keys []string, index int, avalues map[string]interface{}, bvalues []map[string]interface{}) map[string]interface{} {
	if len(keys) < 1 {
		if index < 0 || index >= len(bvalues) {
			return nil
		}
		return bvalues[index]
	} else {
		for _, bv := range bvalues {
			if c.matches(keys, avalues, bv) {
				return bv
			}
		}
	}
	return nil
}

func (c sliceCmp) matches(keys []string, avalues map[string]interface{}, bvalues map[string]interface{}) bool {
	for _, key := range keys {
		if avalues[key] != bvalues[key] {
			return false
		}
	}
	return true
}

func (c sliceCmp) convertToStringMaps(b []interface{}) ([]map[string]interface{}, []map[string]interface{}, error) {
	asrc := make([]map[string]interface{}, 0)
	bsrc := make([]map[string]interface{}, 0)
	for _, av := range c.A {
		amap := make(map[string]interface{})
		err := toFromJson(av, &amap)
		if err != nil {
			return nil, nil, err
		}
		asrc = append(asrc, amap)
	}
	for _, bv := range b {
		bmap := make(map[string]interface{})
		err := toFromJson(bv, &bmap)
		if err != nil {
			return nil, nil, err
		}
		bsrc = append(bsrc, bmap)
	}
	return asrc, bsrc, nil
}
