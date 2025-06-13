package mresult

import (
	"fmt"
	"testing"
)

type MResult[T any] struct {
	val      T
	err      error
	executed bool
}

var _ mresult = MResult[any]{}

func new[T any](t *testing.T, val T) MResult[T] {
	return MResult[T]{
		val:      val,
		executed: true,
	}
}

func newErr[T any](t *testing.T, err error) MResult[T] {
	return MResult[T]{
		err:      err,
		executed: true,
	}
}

type newFunc[T any] func(t *testing.T, val T) MResult[T]
type newErrFunc[T any] func(t *testing.T, err error) MResult[T]

func Generator[T any](t *testing.T) (newFunc[T], newErrFunc[T]) {
	t.Helper()
	return new[T], newErr[T]
}

func (r MResult[T]) IsExecuted(t *testing.T) bool {
	t.Helper()
	return r.executed
}

func (r MResult[T]) IsError(t *testing.T) bool {
	t.Helper()
	return r.err != nil
}

func (r MResult[T]) Val(t *testing.T) T {
	t.Helper()
	requireExecuted(t, r)
	return r.val
}

func (r MResult[T]) Err(t *testing.T) error {
	t.Helper()
	requireExecuted(t, r)
	return r.err
}

func (r MResult[T]) HasVal(t *testing.T) (T, bool) {
	t.Helper()
	return r.val, r.executed
}

func (r MResult[T]) HasErr(t *testing.T) (error, bool) {
	t.Helper()
	return r.err, r.executed
}

func (r MResult[T]) name() string {
	return fmt.Sprintf("MResult[%T]", r.val)
}
