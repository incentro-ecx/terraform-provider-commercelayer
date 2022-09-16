package commercelayer

import (
	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/ladydascalie/currency"
)

var currencyCodeValidation = func(i interface{}, path cty.Path) diag.Diagnostics {
	_, err := currency.Get(i.(string))
	return diagErr(err)
}
