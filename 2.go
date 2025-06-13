package mresult

import (
	"fmt"
	"testing"
)

// MResult2 represents a mock result that contains two values and an error.
// This is useful for mocking functions that return two values and an error,
// such as map lookups, type assertions, or functions that return multiple
// related values.
//
// The zero value for MResult2 represents an unexecuted result.
// Use [Generator2] to create instances of this type.
type MResult2[T1 any, T2 any] struct {
	val1     T1
	val2     T2
	err      error
	executed bool
}

var _ mresult = MResult2[any, any]{}

// new2 creates a new successful MResult2 with the given values.
func new2[T1 any, T2 any](t *testing.T, val1 T1, val2 T2) MResult2[T1, T2] {
	t.Helper()
	return MResult2[T1, T2]{
		val1:     val1,
		val2:     val2,
		executed: true,
	}
}

// newErr2 creates a new MResult2 that represents an error case.
func newErr2[T1 any, T2 any](t *testing.T, err error) MResult2[T1, T2] {
	t.Helper()
	return MResult2[T1, T2]{
		err:      err,
		executed: true,
	}
}

// newFunc2 is a function type that creates successful MResult2 instances with two values.
type newFunc2[T1 any, T2 any] func(t *testing.T, val1 T1, val2 T2) MResult2[T1, T2]

// newErrFunc2 is a function type that creates error MResult2 instances.
type newErrFunc2[T1 any, T2 any] func(t *testing.T, err error) MResult2[T1, T2]

// Generator2 returns generator functions for creating [MResult2] instances.
// It returns two functions: one for successful results with two values,
// and one for error results.
//
// Example usage:
//
//	lookupR, lookupRErr := mresult.Generator2[string, int](t)
//
//	testCases := map[string]struct {
//		lookupResult mresult.MResult2[string, int]
//	}{
//		"found": {
//			lookupResult: lookupR(t, "key1", 42),
//		},
//		"database_error": {
//			lookupResult: lookupRErr(t, errors.New("database connection failed")),
//		},
//	}
func Generator2[T1 any, T2 any](t *testing.T) (newFunc2[T1, T2], newErrFunc2[T1, T2]) {
	t.Helper()
	return new2[T1, T2], newErr2[T1, T2]
}

// IsExecuted reports whether this result has been set up for execution.
// It returns true if the result was created by functions returned from [Generator2].
func (r MResult2[T1, T2]) IsExecuted(t *testing.T) bool {
	t.Helper()
	return r.executed
}

// IsError reports whether this result represents an error case.
// It returns true if the result contains a non-nil error.
func (r MResult2[T1, T2]) IsError(t *testing.T) bool {
	t.Helper()
	return r.err != nil
}

// Val1 returns the first value for this result.
// It calls [testing.T.Fatal] internally if the result has not been executed.
func (r MResult2[T1, T2]) Val1(t *testing.T) T1 {
	t.Helper()
	requireExecuted(t, r)
	return r.val1
}

// Val2 returns the second value for this result.
// It calls [testing.T.Fatal] internally if the result has not been executed.
func (r MResult2[T1, T2]) Val2(t *testing.T) T2 {
	t.Helper()
	requireExecuted(t, r)
	return r.val2
}

// Err returns the error value for this result.
// It calls [testing.T.Fatal] internally if the result has not been executed.
func (r MResult2[T1, T2]) Err(t *testing.T) error {
	t.Helper()
	requireExecuted(t, r)
	return r.err
}

// HasVal returns both values and execution status.
// It returns the first value, second value (which may be the zero values for T1 and T2),
// and a boolean indicating whether the result has been executed.
func (r MResult2[T1, T2]) HasVal(t *testing.T) (T1, T2, bool) {
	t.Helper()
	return r.val1, r.val2, r.executed
}

// HasErr returns the error value and execution status.
// It returns the error (which may be nil) and a boolean indicating
// whether the result has been executed.
func (r MResult2[T1, T2]) HasErr(t *testing.T) (error, bool) {
	t.Helper()
	return r.err, r.executed
}

// name returns the string representation of this result type for debugging.
func (r MResult2[T1, T2]) name() string {
	return fmt.Sprintf("MResult2[%T, %T]", r.val1, r.val2)
}
