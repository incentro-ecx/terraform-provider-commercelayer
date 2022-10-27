package commercelayer

import (
	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/ladydascalie/currency"
	"strings"
)

var currencyCodeValidation = func(i interface{}, path cty.Path) diag.Diagnostics {
	_, err := currency.Get(i.(string))
	return diagErr(err)
}

func getInventoryModelStrategies() []string {
	return []string{
		"no_split",
		"split_shipments",
		"ship_from_primary",
		"ship_from_first_available_or_primary",
	}
}

var inventoryModelStrategyValidation = func(i interface{}, path cty.Path) diag.Diagnostics {
	for _, s := range getInventoryModelStrategies() {
		if s == i.(string) {
			return nil
		}
	}
	return diag.Errorf("Invalid inventory model strategy provided: %s. Must be one of %s",
		i.(string), strings.Join(getInventoryModelStrategies(), ", "))
}
