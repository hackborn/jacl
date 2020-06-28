package jacl

// ------------------------------------------------------------
// SINGLE-CMP

// singleCmp compares a single item to another.
type singleCmp struct {
	A interface{} `json:"a,omitempty"`
}

func (c singleCmp) Cmp(b interface{}) (bool, error) {
	// Handle simple comparisons.
	ans, err := compareBasicTypes(c.A, b)
	if err == nil {
		return ans, nil
	}

	// Handle complex comparisons.
	amap := make(map[string]interface{})
	bmap := make(map[string]interface{})
	err = toFromJson(c.A, &amap, b, &bmap)
	if err != nil {
		return false, err
	}
	for k, av := range amap {
		if Compare(av, bmap[k]) == false {
			return false, nil
		}
	}
	return true, nil
}

func (c singleCmp) FactoryKey() string {
	return singleCmpFactoryKey
}
