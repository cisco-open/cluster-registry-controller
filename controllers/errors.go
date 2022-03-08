// Copyright (c) 2021, and 2022 Cisco and/or its affiliates. All rights reserved.

package controllers

import (
	"runtime"

	"emperror.dev/errors"
)

var (
	ErrInvalidClusterID     = errors.New("invalid cluster ID")
	ErrInvalidSecretContent = errors.New("could not found k8s config in secret")
	ErrInvalidSecret        = errors.New("invalid secret type")
	ErrLocalClusterConflict = errors.New("multiple local clusters are defined")
)

func WrapAsPermanentError(err error) error {
	return &permanentError{err, callers(2)}
}

type permanentError struct {
	error
	*stack
}

func (w *permanentError) Cause() error { return w.error }

// Unwrap provides compatibility for Go 1.13 error chains.
func (w *permanentError) Unwrap() error { return w.error }

func (w *permanentError) IsPermanent() bool {
	return true
}

type (
	StackTrace = errors.StackTrace
	Frame      = errors.Frame
	stack      []uintptr
)

// callers is based on the function with the same name in github.com/pkg/errors,
// but accepts a custom depth (useful to customize the error constructor caller depth).
func callers(depth int) *stack {
	const maxDepth = 32

	var pcs [maxDepth]uintptr

	n := runtime.Callers(2+depth, pcs[:])

	var st stack = pcs[0:n]

	return &st
}

func (s *stack) StackTrace() StackTrace {
	f := make([]Frame, len(*s))
	for i := 0; i < len(f); i++ {
		f[i] = Frame((*s)[i])
	}

	return f
}
