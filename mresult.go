package mresult

import (
	"testing"
)

type mresult interface {
	name() string
	IsExecuted(t *testing.T) bool
	IsError(t *testing.T) bool
	HasErr(t *testing.T) (error, bool)
}

func requireExecuted(t *testing.T, r mresult) {
	t.Helper()
	if !r.IsExecuted(t) {
		t.Fatalf("%s is not expected to execute", r.name())
	}
}
