package jacl

import (
	"encoding/json"
)

// ------------------------------------------------------------
// CMPER-FACTORY

// CmperFactory wraps a Cmper so it can be marshalled to
// and from JSON.
type CmperFactory struct {
	Cmper Cmper `json:"cmper,omitempty"`
}

// Cmp() is a convenience function for doing the comparison.
// It won't return an error (you can directly run the
// comparison for that), and it will return true if there is
// no comparison to run.
func (f CmperFactory) Cmp(b interface{}) bool {
	if f.Cmper == nil {
		return true
	}
	ok, _ := f.Cmper.Cmp(b)
	return ok
}

// MarshalJSON() overrides this struct's marshalling to remove the Fields layer.
func (f CmperFactory) MarshalJSON() ([]byte, error) {
	key := ""
	if f.Cmper != nil {
		key = f.Cmper.FactoryKey()
	}
	glue := cmperFactoryGlue{key, f.Cmper}
	return json.Marshal(glue)
}

// UnmarshalJSON() overrides this struct's unmarshalling to remove the Fields layer.
func (f *CmperFactory) UnmarshalJSON(data []byte) error {
	glue := cmperFactoryGlue{}
	err := json.Unmarshal(data, &glue)
	if err != nil {
		return err
	}
	switch glue.Key {
	case singleCmpFactoryKey:
		c := &singleCmp{}
		err = toFromJson(glue.Cmper, c)
		f.Cmper = c
	case sliceCmpFactoryKey:
		c := &sliceCmp{}
		err = toFromJson(glue.Cmper, c)
		f.Cmper = c
	}
	return err
}

type cmperFactoryGlue struct {
	Key   string      `json:"key,omitempty"`
	Cmper interface{} `json:"cmper,omitempty"`
}

// ------------------------------------------------------------
// CONST and VAR

const (
	singleCmpFactoryKey = "jacl-singlecmp"
	sliceCmpFactoryKey  = "jacl-slicecmp"
)
