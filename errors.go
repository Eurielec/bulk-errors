// Licensed under the GPLv3, see LICENCE file for details.

package bulkerrs

import (
	"bytes"
	"fmt"

	"github.com/juju/errors"
)

// Errs holds an array of JujuErr (juju/errors).
//
// It may be embedded in also custom error types that have been converted to JujuErr.
type Errs struct {
	// holds inner errors
	errors []JujuErr
}

// If an error complies this interface, then it's a juju error :)
type JujuErr interface {
	Cause() error
	Error() string
	Format(fmt.State, rune)
	Location() (string, int)
	Message() string
	SetLocation(int)
	StackTrace() []string
	Underlying() error
}

// Juju Error function interfaces
type NewXFn func(error, string) error

// Aux Juju Error constructor, panics if argument error is nil
func getJujuErr(err error) JujuErr {
	if err == nil {
		panic("Error is nil")
	}
	jErr, ok := err.(JujuErr)
	if !ok {
		newerr := errors.Annotate(err, err.Error())
		jErr = newerr.(JujuErr)
		jErr.SetLocation(2)
	}
	return jErr
}

func NewErr() Errs {
	return Errs{}
}

func NewErrOr(err error) Errs {
	if err == nil {
		return Errs{}
	} else if myerr, ok := err.(*Errs); ok {
		return *myerr
	} else {
		return Errs{
			errors: []JujuErr{getJujuErr(err)},
		}
	}
}

func (e *Errs) NewErr(format string, args ...interface{}) {
	err := errors.NewErr(format, args...)
	err.SetLocation(1)
	e.errors = append(e.errors, &err)
}

func (e *Errs) NewErrWithCause(other error, format string, args ...interface{}) {
	err := errors.NewErrWithCause(other, format, args...)
	err.SetLocation(1)
	e.errors = append(e.errors, &err)
}

func (e *Errs) Append(errs ...error) bool {
	result := false
	for _, err := range errs {
		if err == nil {
			// pass
		} else if myerr, ok := err.(*Errs); ok {
			e.errors = append(e.errors, myerr.errors...)
		} else {
			e.errors = append(e.errors, getJujuErr(err))
		}
		result = true
	}
	return result
}

// Appends error if condition, returns condition
func (e *Errs) AppendIf(condition bool, msg string) bool {
	if condition {
		err := errors.New(msg).(JujuErr)
		err.SetLocation(1)
		e.Append(err)
	}
	return condition
}

func (e *Errs) AppendIfX(condition bool, newErr NewXFn, msg string, other error) bool {
	if condition {
		err := newErr(other, msg).(JujuErr)
		err.SetLocation(1)
		e.Append(err)
	}
	return condition
}

func Concat(errs ...error) Errs {
	err := NewErr()
	for _, erri := range errs {
		err.Append(erri)
	}
	return err
}

func (e Errs) ToError() error {
	if len(e.errors) == 0 {
		return nil
	}
	return &e
}

func (e *Errs) Error() string {
	var buffer bytes.Buffer
	for i, err := range e.errors {
		if i != 0 {
			buffer.WriteString("\n")
		}
		buffer.WriteString(err.Error())
	}
	return buffer.String()
}

func (e *Errs) Errors() []string {
	errs := make([]string, len(e.errors))
	for i, err := range e.errors {
		errs[i] = err.Error()
	}
	return errs
}

func (e *Errs) InnerErrors() []error {
	errs := make([]error, len(e.errors))
	for i, err := range e.errors {
		errs[i] = err
	}
	return errs
}

func (e Errs) Format(s fmt.State, verb rune) {
	fmt.Fprintf(s, "[")
	max_idx := len(e.errors) - 1
	for i, err := range e.errors {
		err.Format(s, verb)
		if i != max_idx {
			fmt.Fprintf(s, ", ")
		}
	}
	fmt.Fprintf(s, "]")
}
