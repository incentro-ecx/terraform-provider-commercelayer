package commercelayer

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCurrencyCodeValidationErr(t *testing.T) {
	diag := currencyCodeValidation("FOOBAR", nil)
	assert.True(t, diag.HasError())
}

func TestCurrencyCodeValidationOK(t *testing.T) {
	diag := currencyCodeValidation("EUR", nil)
	assert.False(t, diag.HasError())
}
