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

func TestPaymentSourceValidationError(t *testing.T) {
	diag := paymentSourceValidation("Adyen", nil)
	assert.True(t, diag.HasError())
}

func TestPaymentSourceValidationOK(t *testing.T) {
	diag := paymentSourceValidation("BraintreePayment", nil)
	assert.False(t, diag.HasError())
}
