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

	// Handle complex comparisons.
	amap := make(map[string]interface{})
	bmap := make(map[string]interface{})
	err = toFromJson(c.A, &amap, b, &bmap)
	if err != nil {
		return newEvaluationError(err)
	}
	for k, av := range amap {
		if Compare(av, bmap[k]) == false {
			return newComparisonError(fmt.Sprintf("have %v want %v", toJson(b), toJson(c.A)))
		}
	}
	return nil
}

func (c singleCmp) SerializeKey() string {
	return singleCmpFactoryKey
}
