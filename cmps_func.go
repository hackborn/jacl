package jacl

import (
	"fmt"
)

type CmpsFunc interface {
	Eval(resp []interface{}) error
}

type sizeis struct {
	Size int `json:"size,omitempty"`
}

func (f sizeis) Eval(resp []interface{}) error {
	if len(resp) == f.Size {
		return nil
	}
	return fmt.Errorf("Size mismatch, have %v want %v", len(resp), f.Size)
}
