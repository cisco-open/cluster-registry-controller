// Copyright (c) 2021, and 2022 Cisco and/or its affiliates. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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
