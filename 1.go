package mresult

import (
	"fmt"
	"testing"
)

// MResult represents a mock result that contains a single value and an error.
// This is useful for mocking functions that return a value and an error, such as
// database queries, API calls, or file operations.
//
// The zero value for MResult represents an unexecuted result.
// Use [Generator] to create instances of this type.
type MResult[T any] struct {
	val      T
	err      error
	executed bool
}

var _ mresult = MResult[any]{}

// new creates a new successful MResult with the given value.
func new[T any](t *testing.T, val T) MResult[T] {
	return MResult[T]{
		val:      val,
		executed: true,
	}
}

// newErr creates a new MResult that represents an error case.
func newErr[T any](t *testing.T, err error) MResult[T] {
	return MResult[T]{
		err:      err,
		executed: true,
	}
}

// newFunc is a function type that creates successful MResult instances with a value.
type newFunc[T any] func(t *testing.T, val T) MResult[T]

// newErrFunc is a function type that creates error MResult instances.
type newErrFunc[T any] func(t *testing.T, err error) MResult[T]

// Generator returns generator functions for creating [MResult] instances.
// It returns two functions: one for successful results with a value,
// and one for error results.
//
// Example usage:
//
//	getUserR, getUserRErr := mresult.Generator[*User](t)
//
//	testCases := map[string]struct {
//		getUserResult mresult.MResult[*User]
//	}{
//		"success": {
//			getUserResult: getUserR(t, &User{ID: 1, Name: "John"}),
//		},
//		"not_found": {
//			getUserResult: getUserRErr(t, errors.New("user not found")),
//		},
//	}
func Generator[T any](t *testing.T) (newFunc[T], newErrFunc[T]) {
	t.Helper()
	return new[T], newErr[T]
}

// IsExecuted reports whether this result has been set up for execution.
// It returns true if the result was created by functions returned from [Generator].
func (r MResult[T]) IsExecuted(t *testing.T) bool {
	t.Helper()
	return r.executed
}

// IsError reports whether this result represents an error case.
// It returns true if the result contains a non-nil error.
func (r MResult[T]) IsError(t *testing.T) bool {
	t.Helper()
	return r.err != nil
}

// Val returns the value for this result.
// It calls [testing.T.Fatal] internally if the result has not been executed.
func (r MResult[T]) Val(t *testing.T) T {
	t.Helper()
	requireExecuted(t, r)
	return r.val
}

// Err returns the error value for this result.
// It calls [testing.T.Fatal] internally if the result has not been executed.
func (r MResult[T]) Err(t *testing.T) error {
	t.Helper()
	requireExecuted(t, r)
	return r.err
}

// HasVal returns the value and execution status.
// It returns the value (which may be the zero value for T) and a boolean
// indicating whether the result has been executed.
func (r MResult[T]) HasVal(t *testing.T) (T, bool) {
	t.Helper()
	return r.val, r.executed
}

// HasErr returns the error value and execution status.
// It returns the error (which may be nil) and a boolean indicating
// whether the result has been executed.
func (r MResult[T]) HasErr(t *testing.T) (error, bool) {
	t.Helper()
	return r.err, r.executed
}

// name returns the string representation of this result type for debugging.
func (r MResult[T]) name() string {
	return fmt.Sprintf("MResult[%T]", r.val)
}
