// The MIT License (MIT)

// Copyright (c) 2016 Meteora

// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:

// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.

// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.
package errors

import "sync"

// ErrorList is used to chain a list of potential errors and is thread-safe.
type ErrorList struct {
	mux  sync.RWMutex
	errs []error
}

// Error will return the string-form of the errors.
// Implements the error interface.
func (e *ErrorList) Error() string {
	if e == nil {
		return ""
	}

	e.mux.RLock()
	defer e.mux.RUnlock()

	if len(e.errs) == 0 {
		return ""
	}

	if len(e.errs) == 1 {
		return e.errs[0].Error()
	}

	b := []byte("the following errors occured:\n")
	for _, err := range e.errs {
		b = append(b, err.Error()...)
		b = append(b, '\n')
	}

	return string(b)
}

// Err will return an error if the errorlist is not empty.
// If there's only 1 error, it will be directly returned.
// If the errorlist is empty - nil is returned.
func (e *ErrorList) Err() (err error) {
	if e == nil {
		return
	}
	e.mux.RLock()
	switch len(e.errs) {
	case 0: // do nothing
	case 1:
		err = e.errs[0]
	default:
		err = e
	}
	e.mux.RUnlock()
	return
}

// Push will push an error to the errorlist
// If err is a errorlist, it will be merged.
// If the errorlist is nil, it will be created.
func (e *ErrorList) Push(err error) {
	if err == nil {
		return
	}

	e.mux.Lock()
	defer e.mux.Unlock()

	switch v := err.(type) {
	case *ErrorList:
		v.ForEach(func(err error) {
			e.errs = append(e.errs, err)
		})

	default:
		e.errs = append(e.errs, err)
	}
}

// ForEach will iterate through all of the errors within the error list.
func (e *ErrorList) ForEach(fn func(error)) {
	if e == nil {
		return
	}

	e.mux.RLock()
	for _, err := range e.errs {
		fn(err)
	}
	e.mux.RUnlock()
}

// Copy will copy the items from the inbound error list to the source
func (e *ErrorList) Copy(in *ErrorList) {
	if in == nil {
		return
	}

	e.mux.Lock()
	defer e.mux.Unlock()

	in.mux.RLock()
	defer in.mux.RUnlock()

	e.errs = append(e.errs, in.errs...)
}

// Len will return the length of the inner errors list.
func (e *ErrorList) Len() (n int) {
	if e == nil {
		return
	}

	e.mux.RLock()
	n = len(e.errs)
	e.mux.RUnlock()
	return
}
