package jacl

// ------------------------------------------------------------
// SLICE-CMP

// sliceCmp compares a list of items to another list.
// It has an optional key: If present, the key is used
// as an identifier to match items in each list. If absent,
// it assumes the list are in the correct order.
type sliceCmp struct {
	Key string        `json:"key,omitempty"`
	A   []interface{} `json:"a,omitempty"`
}

func (c sliceCmp) Cmp(b interface{}) (bool, error) {
	if c.A == nil && b == nil {
		return true, nil
	}
	bslice := make([]interface{}, 0, 0)
	err := toFromJson(b, &bslice)
	if err != nil {
		return false, err
	}

	asrc, bsrc, err := c.convertToStringMaps(bslice)
	if err == nil {
		return c.cmpStringMaps(asrc, bsrc)
	}

	// If I couldn't convert to string maps, assume the slices
	// contain literals.
	return c.cmpSlices(c.A, bslice)
}

func (c sliceCmp) FactoryKey() string {
	return sliceCmpFactoryKey
}

func (c sliceCmp) cmpStringMaps(asrc, bsrc []map[string]interface{}) (bool, error) {
	for i, av := range asrc {
		bv := c.find(c.Key, i, av, bsrc)
		if bv == nil {
			return false, nil
		}
		if Compare(av, bv) == false {
			return false, nil
		}
	}
	return true, nil
}

func (c sliceCmp) cmpSlices(aslice, bslice []interface{}) (bool, error) {
	if len(aslice) != len(bslice) {
		return false, nil
	}
	for i, av := range aslice {
		if Compare(av, bslice[i]) == false {
			return false, nil
		}
	}
	return true, nil
}

func (c sliceCmp) find(key string, index int, avalues map[string]interface{}, bvalues []map[string]interface{}) map[string]interface{} {
	if key == "" {
		if index < 0 || index >= len(bvalues) {
			return nil
		}
		return bvalues[index]
	} else {
		if k, ok := avalues[key]; k != nil && ok {
			for _, bv := range bvalues {
				if bv[key] == k {
					return bv
				}
			}
		}
	}
	return nil
}

func (c sliceCmp) convertToSlices(b []interface{}) ([]interface{}, []interface{}, error) {
	asrc := make([]interface{}, 0, 0)
	bsrc := make([]interface{}, 0, 0)
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

func (c sliceCmp) convertToStringMaps(b []interface{}) ([]map[string]interface{}, []map[string]interface{}, error) {
	asrc := make([]map[string]interface{}, 0, 0)
	bsrc := make([]map[string]interface{}, 0, 0)
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
