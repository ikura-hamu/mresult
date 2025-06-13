package mresult

import (
	"fmt"
	"testing"
)

type MResult2[T1 any, T2 any] struct {
	val1     T1
	val2     T2
	err      error
	executed bool
}

var _ mresult = MResult2[any, any]{}

func new2[T1 any, T2 any](t *testing.T, val1 T1, val2 T2) MResult2[T1, T2] {
	t.Helper()
	return MResult2[T1, T2]{
		val1:     val1,
		val2:     val2,
		executed: true,
	}
}

func newErr2[T1 any, T2 any](t *testing.T, err error) MResult2[T1, T2] {
	t.Helper()
	return MResult2[T1, T2]{
		err:      err,
		executed: true,
	}
}

type newFunc2[T1 any, T2 any] func(t *testing.T, val1 T1, val2 T2) MResult2[T1, T2]
type newErrFunc2[T1 any, T2 any] func(t *testing.T, err error) MResult2[T1, T2]

func Generator2[T1 any, T2 any](t *testing.T) (newFunc2[T1, T2], newErrFunc2[T1, T2]) {
	t.Helper()
	return new2[T1, T2], newErr2[T1, T2]
}

func (r MResult2[T1, T2]) IsExecuted(t *testing.T) bool {
	t.Helper()
	return r.executed
}

func (r MResult2[T1, T2]) IsError(t *testing.T) bool {
	t.Helper()
	return r.err != nil
}

func (r MResult2[T1, T2]) Val1(t *testing.T) T1 {
	t.Helper()
	requireExecuted(t, r)
	return r.val1
}

func (r MResult2[T1, T2]) Val2(t *testing.T) T2 {
	t.Helper()
	requireExecuted(t, r)
	return r.val2
}

func (r MResult2[T1, T2]) Err(t *testing.T) error {
	t.Helper()
	requireExecuted(t, r)
	return r.err
}

func (r MResult2[T1, T2]) HasVal(t *testing.T) (T1, T2, bool) {
	t.Helper()
	return r.val1, r.val2, r.executed
}

func (r MResult2[T1, T2]) HasErr(t *testing.T) (error, bool) {
	t.Helper()
	return r.err, r.executed
}

func (r MResult2[T1, T2]) name() string {
	return fmt.Sprintf("MResult2[%T, %T]", r.val1, r.val2)
}
