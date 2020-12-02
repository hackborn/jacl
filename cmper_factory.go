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

// Cmp is a convenience function for doing the comparison.
func (f CmperFactory) Cmp(b interface{}) error {
	if f.Cmper == nil {
		return nil
	}
	return f.Cmper.Cmp(b)
}

// MarshalJSON overrides this struct's marshalling to remove the Fields layer.
func (f CmperFactory) MarshalJSON() ([]byte, error) {
	key := ""
	if f.Cmper != nil {
		if s, ok := f.Cmper.(serializer); ok {
			key = s.SerializeKey()
		}
	}
	glue := cmperFactoryGlue{key, f.Cmper}
	return json.Marshal(glue)
}

// UnmarshalJSON overrides this struct's unmarshalling to remove the Fields layer.
func (f *CmperFactory) UnmarshalJSON(data []byte) error {
	glue := cmperFactoryGlue{}
	err := json.Unmarshal(data, &glue)
	if err != nil {
		return err
	}
	switch glue.Key {
	case nilCmpFactoryKey:
		f.Cmper = &nilCmp{}
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
// SERIALIZER

// serializer defines items that can be serialized to a CmperFactory.
type serializer interface {
	// Answer a unique key so I can be reinstantiated after marshalling.
	SerializeKey() string
}

// ------------------------------------------------------------
// CONST and VAR

const (
	nilCmpFactoryKey    = "jacl-nilcmp"
	singleCmpFactoryKey = "jacl-singlecmp"
	sliceCmpFactoryKey  = "jacl-slicecmp"
)
