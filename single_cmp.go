package jacl

import (
	"fmt"
)

// ------------------------------------------------------------
// SINGLE-CMP

// singleCmp compares a single item to another.
type singleCmp struct {
	A interface{} `json:"a,omitempty"`
}

func (c singleCmp) Cmp(b interface{}) error {
	// Handle simple comparisons.
	ans, err := compareBasicTypes(c.A, b)
	if err == nil {
		if ans {
			return nil
		}
		return newComparisonError(fmt.Sprintf("have %v want %v", toJson(b), toJson(c.A)))
	}

	// Handle slice comparisons.
	handled, err := c.cmpAsSlices(c.A, b)
	if handled {
		return err
	}

	// Handle map comparisons.
	amap := make(map[string]interface{})
	bmap := make(map[string]interface{})
	err = toFromJson(c.A, &amap, b, &bmap)
	if err != nil {
		return newEvaluationError(err)
	}
	for k, av := range amap {
		if !compare(av, bmap[k]) {
			return newComparisonError(fmt.Sprintf("have %v want %v", toJson(b), toJson(c.A)))
		}
	}
	return nil
}

func (c singleCmp) SerializeKey() string {
	return singleCmpFactoryKey
}

func (c singleCmp) cmpAsSlices(_a, _b interface{}) (bool, error) {
	// Need to encode/decode the data to eliminate variations in slice type
	var aslice []interface{}
	var bslice []interface{}
	err := toFromJson(_a, &aslice, _b, &bslice)
	if err != nil {
		return false, err
	}
	if len(aslice) != len(bslice) {
		return true, newComparisonError(fmt.Sprintf("have length %v want length %v", len(bslice), len(aslice)))
	}
	for i, ai := range aslice {
		if !compare(ai, bslice[i]) {
			return true, newComparisonError(fmt.Sprintf("have %v want %v", toJson(bslice), toJson(aslice)))
		}
	}
	return true, nil
}
