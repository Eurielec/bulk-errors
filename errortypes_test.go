// Licensed under the GPLv3, see LICENCE file for details.

package bulkerrs_test

import (
	"testing"

	"github.com/eurielec/bulkerrs"
	"github.com/juju/errors"
	"github.com/stretchr/testify/assert"
)

var isXFns = [...]bulkerrs.IsXFn{
	bulkerrs.IsTimeout,
	bulkerrs.IsNotFound,
	bulkerrs.IsUserNotFound,
	bulkerrs.IsUnauthorized,
	bulkerrs.IsNotImplemented,
	bulkerrs.IsAlreadyExists,
	bulkerrs.IsNotSupported,
	bulkerrs.IsNotValid,
	bulkerrs.IsNotProvisioned,
	bulkerrs.IsNotAssigned,
	bulkerrs.IsBadRequest,
	bulkerrs.IsMethodNotAllowed,
	bulkerrs.IsForbidden,
}

var newXFns = [...]bulkerrs.NewXFn{
	errors.NewTimeout,
	errors.NewNotFound,
	errors.NewUserNotFound,
	errors.NewUnauthorized,
	errors.NewNotImplemented,
	errors.NewAlreadyExists,
	errors.NewNotSupported,
	errors.NewNotValid,
	errors.NewNotProvisioned,
	errors.NewNotAssigned,
	errors.NewBadRequest,
	errors.NewMethodNotAllowed,
	errors.NewForbidden,
}

func TestBase(t *testing.T) {
	// Check errors type of a nil
	var errs *bulkerrs.Errs
	for _, isx := range isXFns {
		assert.False(t, isx(errs), "Should return false")
	}

	// Check errors type of a JujuErr
	for i, newx := range newXFns {
		for j, isx := range isXFns {
			assert.Equal(t, i == j, isx(newx(nil, "")), "Should return only one true in loop")
		}
	}
}

func TestX(t *testing.T) {
	for i := 0; i < len(isXFns); i++ {
		// New Errs
		errs := bulkerrs.NewErr()

		// Fill with errors
		for j, newx := range newXFns {
			errs.AppendIfX(i == j, newx, "", nil)
		}

		// New Errs -> error
		err := errs.ToError()

		// Check errors type
		for j, isx := range isXFns {
			assert.Equal(t, i == j, isx(err), "Should return only one true in loop")
		}

		// New Errs
		errs = bulkerrs.NewErr()

		// Fill with error
		for j, newx := range newXFns {
			errs.AppendIfX(i != j, newx, "", nil)
		}

		// New Errs -> error
		err = errs.ToError()

		// Check errors type
		for j, isx := range isXFns {
			assert.Equal(t, i != j, isx(err), "Should return only one false in loop")
		}
	}
}
