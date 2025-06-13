package mresult

import (
	"testing"
)

type MResult0 struct {
	err      error
	executed bool
}

var _ mresult = MResult0{}

func new0(t *testing.T) MResult0 {
	t.Helper()
	return MResult0{
		executed: true,
	}
}

func newErr0(t *testing.T, err error) MResult0 {
	t.Helper()
	return MResult0{
		err:      err,
		executed: true,
	}
}

type newFunc0 func(t *testing.T) MResult0
type newErrFunc0 func(t *testing.T, err error) MResult0

func Generator0(t *testing.T) (newFunc0, newErrFunc0) {
	t.Helper()
	return new0, newErr0
}

func (r MResult0) IsExecuted(t *testing.T) bool {
	t.Helper()
	return r.executed
}

func (r MResult0) IsError(t *testing.T) bool {
	t.Helper()
	return r.err != nil
}

func (r MResult0) Err(t *testing.T) error {
	t.Helper()
	requireExecuted(t, r)
	return r.err
}

func (r MResult0) HasErr(t *testing.T) (error, bool) {
	t.Helper()
	return r.err, r.executed
}

func (r MResult0) name() string {
	return "MResult0"
}
