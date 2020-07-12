package jacl

import (
	"encoding/json"
	"fmt"
)

// ------------------------------------------------------------
// CMPS-FUNC

// CmpsFunc defines behaviour objects that can be sent to a Cmps().
type CmpsFunc interface {
	Eval(resp []interface{}) error

	// Answer a unique key so I can be reinstantiated after marshalling.
	FactoryKey() string
}

// ------------------------------------------------------------
// KEY-FN FUNCTION

// keyFn defines what keys to use when identifying an item in a slice.
type keyFn struct {
	Keys []string `json:"keys,omitempty"`
}

func (f keyFn) Eval(resp []interface{}) error {
	return nil
}

func (f keyFn) FactoryKey() string {
	return keyFactoryKey
}

// ------------------------------------------------------------
// SIZEIS-FN FUNCTION

// sizeis is a function to evaluate the size of a slice.
type sizeisFn struct {
	Size int `json:"size,omitempty"`
}

func (f sizeisFn) Eval(resp []interface{}) error {
	if len(resp) == f.Size {
		return nil
	}
	return newComparisonError(fmt.Sprintf("Size mismatch, have %v want %v", len(resp), f.Size))
}

func (f sizeisFn) FactoryKey() string {
	return sizeisFactoryKey
}

// ------------------------------------------------------------
// FUNC-FACTORY

// FuncFactory wraps a CmpsFunc so it can be marshalled to
// and from JSON.
type FuncFactory struct {
	Fn CmpsFunc `json:"func,omitempty"`
}

// Cmp() is a convenience function for doing the comparison.
// It won't return an error (you can directly run the
// comparison for that), and it will return true if there is
// no comparison to run.
func (f FuncFactory) Eval(resp []interface{}) error {
	if f.Fn == nil {
		return nil
	}
	return f.Fn.Eval(resp)
}

// MarshalJSON() overrides this struct's marshalling to remove the Fields layer.
func (f FuncFactory) MarshalJSON() ([]byte, error) {
	key := ""
	if f.Fn != nil {
		key = f.Fn.FactoryKey()
	}
	glue := funcFactoryGlue{key, f.Fn}
	return json.Marshal(glue)
}

// UnmarshalJSON() overrides this struct's unmarshalling to remove the Fields layer.
func (f *FuncFactory) UnmarshalJSON(data []byte) error {
	glue := funcFactoryGlue{}
	err := json.Unmarshal(data, &glue)
	if err != nil {
		return err
	}
	switch glue.Key {
	case keyFactoryKey:
		fn := &keyFn{}
		err = toFromJson(glue.Fn, fn)
		f.Fn = fn
	case sizeisFactoryKey:
		fn := &sizeisFn{}
		err = toFromJson(glue.Fn, fn)
		f.Fn = fn
	}
	return err
}

type funcFactoryGlue struct {
	Key string      `json:"key,omitempty"`
	Fn  interface{} `json:"fn,omitempty"`
}

// ------------------------------------------------------------
// CONST and VAR

const (
	keyFactoryKey    = "jacl-key"
	sizeisFactoryKey = "jacl-sizeis"
)
