package support

import (
	"errors"
	"runtime"
)

var InvalidCollectionKeyError = errors.New("")
var CanNotCreateMapError error = errors.New("can not create map")
var CanNotFoundValueError = errors.New("can not found value")
var CanNotInstantiateCallbackWithParameters = errors.New("Can not instantiate callback with parameters")

type withStack struct {
	error
	*stack
}

// stack represents a stack of program counters.
type stack []uintptr

func WithStack(err error) error {
	if err == nil {
		return nil
	}
	return &withStack{
		err,
		callers(),
	}
}

func callers() *stack {
	const depth = 32
	var pcs [depth]uintptr
	n := runtime.Callers(3, pcs[:])
	var st stack = pcs[0:n]
	return &st
}
