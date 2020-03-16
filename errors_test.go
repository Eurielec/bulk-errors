// Licensed under the GPLv3, see LICENCE file for details.

package bulkerrs_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/eurielec/bulkerrs"
	"github.com/juju/errors"
	"github.com/stretchr/testify/assert"
)

// Example errors
var /* const */ err_std error = fmt.Errorf("Opa!")
var /* const */ err_std_null error = nil

func TestConstructor(t *testing.T) {
	// New Errs
	errs := bulkerrs.NewErr()
	assert.Equal(t, 0, len(errs.InnerErrors()), "Should return an Errs without errors")

	// ValidOrNil
	assert.Nil(t, errs.ValidOrNil(), "Should return a nil")

	// Get a Copy
	assert.Equal(t, errs, bulkerrs.NewErrOr(&errs), "Should return an Errs copy")

	// Get or create Errs
	assert.Equal(t, bulkerrs.Errs{}, bulkerrs.NewErrOr(err_std_null), "Should return a brand new Errs")

	// Get or create Errs, but with valid error
	errs = bulkerrs.NewErrOr(err_std)
	assert.Equal(t, 1, len(errs.InnerErrors()), "Should return a new Errs with the error appended")

	// ValidOrNil
	errs_ptr, ok := errs.ValidOrNil().(*bulkerrs.Errs)
	assert.True(t, ok, "Should cast")
	assert.Equal(t, errs, *errs_ptr, "ValidOrNil should return a pointer to Errs")

	// Get a Copy
	assert.Equal(t, errs, bulkerrs.NewErrOr(errs.ValidOrNil()), "Should return an Errs copy")
}

func TestAppend(t *testing.T) {
	// Init
	errs := bulkerrs.NewErr()

	// Append error
	errs.NewErr("Oopsie!")
	assert.Equal(t, 1, len(errs.InnerErrors()), "Errs should have 1 error")

	// Append wrapped error
	errs.NewErrWithCause(err_std, "Let me tell you a little more...")
	assert.Equal(t, 2, len(errs.InnerErrors()), "Errs should have 2 errors")

	// Attempt to append a null error
	errs.Append(err_std_null)
	assert.Equal(t, 2, len(errs.InnerErrors()), "Errs should still have 2 errors")

	// Append a valid error
	errs.Append(err_std)
	assert.Equal(t, 3, len(errs.InnerErrors()), "Errs should have 3 errors")

	// Append an Errs
	errs_aux := bulkerrs.NewErrOr(errs.InnerErrors()[2])
	errs.Append(errs_aux.ValidOrNil())
	assert.Equal(t, 4, len(errs.InnerErrors()), "Errs should have 4 errors")
	assert.Equal(t, errs.InnerErrors()[2], errs.InnerErrors()[3], "Errs should be the same error")

	// AppendIf False
	condition := errs.AppendIf(false, "This will never fail")
	assert.Equal(t, 4, len(errs.InnerErrors()), "Errs should still have 4 errors")
	assert.False(t, condition, "Should return false")

	// AppendIf True
	condition = errs.AppendIf(true, "This failed!")
	assert.Equal(t, 5, len(errs.InnerErrors()), "Errs should have 5 errors")
	assert.True(t, condition, "Should return true")

	// AppendIfX True
	condition = errs.AppendIfX(true, errors.NewNotValid, "This specially failed!", err_std)
	assert.Equal(t, 6, len(errs.InnerErrors()), "Errs should have 6 errors")
	assert.Equal(t, err_std, errs.InnerErrors()[5].(bulkerrs.JujuErr).Underlying(), "Errs[5] should contain err_std as previous error")
	assert.True(t, condition, "Should return true")
}

func TestAuxFuncts(t *testing.T) {
	// Init
	errs1 := bulkerrs.NewErr()

	// Append a few errors
	errs1.NewErr("Oopsie!")
	errs1.NewErrWithCause(err_std, "Let me tell you a little more...")
	errs1.Append(err_std)

	// Concat errors
	errs_out := bulkerrs.Concat(errs1.ValidOrNil(), errs1.ValidOrNil())
	errs1_errs := errs1.InnerErrors()
	len_errs1_errs := len(errs1_errs)

	// Get InnerErrors
	for i, err := range errs_out.InnerErrors() {
		assert.Equal(t, errs1_errs[i%len_errs1_errs], err, "Errs_out.InnerErrors() should contain Errs1 errors twice")
	}

	// Get .Error() slice
	for i, err_str := range errs_out.Errors() {
		assert.Equal(t, errs1_errs[i%len_errs1_errs].Error(), err_str, "Errs_out.Errors() should contain Errs1 errors.Error() twice")
	}

	// Get Errs.Error()
	assert.Equal(t, strings.Join(errs_out.Errors(), "\n"), errs_out.Error(), "Errs_out.Error() should be a join of Errs1 errors.Error()")

	// Init2
	errs2 := bulkerrs.NewErr()
	errs2.NewErr("Oopsie!")
	errs2 = bulkerrs.Concat(errs2.ValidOrNil(), err_std)

	// Get fmt.Format()
	for _, verb := range []string{"s", "+v", "#v", "q"} {
		inner_errs2 := errs2.InnerErrors()
		errs_fmt := make([]string, len(inner_errs2))
		for i, err := range inner_errs2 {
			errs_fmt[i] = fmt.Sprintf("%"+fmt.Sprintf("%s", verb), err)
		}
		err_fmt1 := fmt.Sprintf("[" + strings.Join(errs_fmt, ", ") + "]")
		err_fmt2 := fmt.Sprintf("%"+fmt.Sprintf("%s", verb), errs2)
		assert.Equal(t, err_fmt1, err_fmt2, "Should be a join of Errs2 errors.Format()")
	}
}
