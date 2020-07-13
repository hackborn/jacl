package jacl

import (
	"errors"
)

// ------------------------------------------------------------
// COMPARISON-ERROR

// ComparisonError indicates that a comparison failed.
type ComparisonError struct {
	s string
}

func newComparisonError(s string) error {
	return &ComparisonError{s: s}
}

func (e *ComparisonError) Error() string {
	return e.s
}

// ------------------------------------------------------------
// EVALUATION-ERROR

// EvaluationError indicates a comparison could not be performed.
type EvaluationError struct {
	err error
}

func newEvaluationError(err error) error {
	return &EvaluationError{err: err}
}

func (e *EvaluationError) Error() string {
	return e.err.Error()
}

func (e *EvaluationError) Unwrap() error {
	return e.err
}

// ------------------------------------------------------------
// TESTING

func equalErr(a, b error) bool {
	var ce *ComparisonError
	var ee *EvaluationError

	if a == b {
		return true
	} else if a == nil {
		return false
	} else if b == nil {
		return false
	} else if errors.As(a, &ce) && errors.As(b, &ce) {
		return true
	} else if errors.As(a, &ee) && errors.As(b, &ee) {
		return true
	} else {
		return a.Error() == b.Error()
	}
}

// ------------------------------------------------------------
// CONST and VAR

const (
	haveWantFmt       = "have %v want %v"
	haveWantLengthFmt = "have length %v want length %v"
)
