// Package mresult provides utilities for simplifying mock return values in Go tests.
//
// This package helps manage mock results in a clean and type-safe way, reducing
// boilerplate code when writing table-driven tests with mocks. It supports various
// return value configurations including single values, multiple values, and error cases.
//
// The main types are [MResult], [MResult0], and [MResult2] for handling different
// numbers of return values, and their corresponding generator functions [Generator],
// [Generator0], and [Generator2].
//
// Example usage:
//
//	// Create generators for mock results
//	getUserR, getUserRErr := mresult.Generator[*User](t)
//
//	// Use in test cases
//	testCases := map[string]struct {
//		getUserResult mresult.MResult[*User]
//	}{
//		"success": {
//			getUserResult: getUserR(t, &User{ID: 1}),
//		},
//		"error": {
//			getUserResult: getUserRErr(t, errors.New("database error")),
//		},
//	}
package mresult

import (
	"testing"
)

// mresult is the internal interface that defines common behavior for all mock result types.
// It provides methods to check execution status, error status, and retrieve error values.
type mresult interface {
	name() string
	IsExecuted(t *testing.T) bool
	IsError(t *testing.T) bool
	HasErr(t *testing.T) (error, bool)
}

// requireExecuted verifies that a mock result has been executed.
// It calls [testing.T.Fatalf] if the result has not been executed, helping to catch
// programming errors where uninitialized results are accessed.
func requireExecuted(t *testing.T, r mresult) {
	t.Helper()
	if !r.IsExecuted(t) {
		t.Fatalf("%s is not expected to execute", r.name())
	}
}
