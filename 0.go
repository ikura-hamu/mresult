package mresult

import (
	"testing"
)

// MResult0 represents a mock result that contains only an error value.
// This is useful for mocking functions that return only an error.
//
// The zero value for MResult0 represents an unexecuted result.
// Use [Generator0] to create instances of this type.
type MResult0 struct {
	err      error
	executed bool
}

var _ mresult = MResult0{}

// new0 creates a new successful MResult0 with no error.
func new0(t *testing.T) MResult0 {
	t.Helper()
	return MResult0{
		executed: true,
	}
}

// newErr0 creates a new MResult0 that represents an error case.
func newErr0(t *testing.T, err error) MResult0 {
	t.Helper()
	return MResult0{
		err:      err,
		executed: true,
	}
}

// newFunc0 is a function type that creates successful MResult0 instances.
type newFunc0 func(t *testing.T) MResult0

// newErrFunc0 is a function type that creates error MResult0 instances.
type newErrFunc0 func(t *testing.T, err error) MResult0

// Generator0 returns generator functions for creating [MResult0] instances.
// It returns two functions: one for successful results and one for error results.
//
// Example usage:
//
//	saveR, saveRErr := mresult.Generator0(t)
//
//	testCases := map[string]struct {
//		saveResult mresult.MResult0
//	}{
//		"success": {
//			saveResult: saveR(t),
//		},
//		"error": {
//			saveResult: saveRErr(t, errors.New("save failed")),
//		},
//	}
func Generator0(t *testing.T) (newFunc0, newErrFunc0) {
	t.Helper()
	return new0, newErr0
}

// IsExecuted reports whether this result has been set up for execution.
// It returns true if the result was created by functions returned from [Generator0].
func (r MResult0) IsExecuted(t *testing.T) bool {
	t.Helper()
	return r.executed
}

// IsError reports whether this result represents an error case.
// It returns true if the result contains a non-nil error.
func (r MResult0) IsError(t *testing.T) bool {
	t.Helper()
	return r.err != nil
}

// Err returns the error value for this result.
// It calls [testing.T.Fatal] internally if the result has not been executed.
func (r MResult0) Err(t *testing.T) error {
	t.Helper()
	requireExecuted(t, r)
	return r.err
}

// HasErr returns the error value and execution status.
// It returns the error (which may be nil) and a boolean indicating
// whether the result has been executed.
func (r MResult0) HasErr(t *testing.T) (error, bool) {
	t.Helper()
	return r.err, r.executed
}

// name returns the string representation of this result type for debugging.
func (r MResult0) name() string {
	return "MResult0"
}
