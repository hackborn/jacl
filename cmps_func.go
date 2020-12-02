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
// NOT-EXISTS-FN FUNCTION

// notExistsFn provides a path to a field that should not exist.
type notExistsFn struct {
	Path []string `json:"path,omitempty"`
}

func (f notExistsFn) Eval(resps []interface{}) error {
	// This is slightly confusing but the functions are designed
	// against slices, so we need to evaluate against each item
	// in the slice.
	for _, resp := range resps {
		if f.existsI(f.Path, resp, false) {
			return newComparisonError(fmt.Sprintf("exists: %v", f.Path))
		}
	}
	return nil
}

func (f notExistsFn) FactoryKey() string {
	return notExistsFactoryKey
}

func (f notExistsFn) existsI(needle []string, _haystack interface{}, converted bool) bool {
	if len(needle) < 1 {
		return false
	}
	switch haystack := _haystack.(type) {
	case string:
		return needle[0] == haystack
	case map[string]interface{}:
		if v, ok := haystack[needle[0]]; ok {
			if len(needle) == 1 {
				return true
			}
			return f.existsI(needle[1:], v, converted)
		} else {
			return false
		}
	default:
		// Convert unknown types into a known format
		if converted {
			panic("unhandled type")
		}
		m := make(map[string]interface{})
		err := toFromJson(_haystack, &m)
		if err != nil {
			panic(err)
		}
		return f.existsI(needle, m, true)
	}
	return false
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
	case notExistsFactoryKey:
		fn := &notExistsFn{}
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
	keyFactoryKey       = "jacl-key"
	notExistsFactoryKey = "jacl-notexists"
	sizeisFactoryKey    = "jacl-sizeis"
)
