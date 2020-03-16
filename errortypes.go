// Licensed under the GPLv3, see LICENCE file for details.

package bulkerrs

import (
	"github.com/juju/errors"
)

type IsXFn func(err error) bool

func (e *Errs) isXloop(isx IsXFn) bool {
	if e == nil {
		return false
	}
	for _, err := range e.errors {
		if isx(err) {
			return true
		}
	}
	return false
}

func isX(isx IsXFn, err error) bool {
	if errs, ok := err.(*Errs); ok {
		return errs.isXloop(isx)
	}
	return isx(err)
}

func IsTimeout(err error) bool {
	return isX(errors.IsTimeout, err)
}

func IsNotFound(err error) bool {
	return isX(errors.IsNotFound, err)
}

func IsUserNotFound(err error) bool {
	return isX(errors.IsUserNotFound, err)
}

func IsUnauthorized(err error) bool {
	return isX(errors.IsUnauthorized, err)
}

func IsNotImplemented(err error) bool {
	return isX(errors.IsNotImplemented, err)
}

func IsAlreadyExists(err error) bool {
	return isX(errors.IsAlreadyExists, err)
}

func IsNotSupported(err error) bool {
	return isX(errors.IsNotSupported, err)
}

func IsNotValid(err error) bool {
	return isX(errors.IsNotValid, err)
}

func IsNotProvisioned(err error) bool {
	return isX(errors.IsNotProvisioned, err)
}

func IsNotAssigned(err error) bool {
	return isX(errors.IsNotAssigned, err)
}

func IsBadRequest(err error) bool {
	return isX(errors.IsBadRequest, err)
}

func IsMethodNotAllowed(err error) bool {
	return isX(errors.IsMethodNotAllowed, err)
}

func IsForbidden(err error) bool {
	return isX(errors.IsForbidden, err)
}
