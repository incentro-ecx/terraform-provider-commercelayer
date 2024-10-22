package commercelayer

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	commercelayer "github.com/incentro-dc/go-commercelayer-sdk/api"
	"net/http"
)

func testAccCheckShippingMethodDestroy(s *terraform.State) error {
	client := testAccProviderCommercelayer.Meta().(*commercelayer.APIClient)

	for _, rs := range s.RootModule().Resources {
		if rs.Type == "commercelayer_shipping_method" {
			err := retryRemoval(10, func() (*http.Response, error) {
				_, resp, err := client.ShippingMethodsApi.GETShippingMethodsShippingMethodId(context.Background(), rs.Primary.ID).
					Execute()
				return resp, err
			})
			if err != nil {
				return err
			}
		}

	}
	return nil
}

func (s *AcceptanceSuite) TestAccShippingMethod_basic() {
	resourceName := "commercelayer_shipping_method.incentro_shipping_method"

	resource.Test(s.T(), resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(s)
		},
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckShippingMethodDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccShippingMethodCreate(resourceName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "type", shippingMethodType),
					resource.TestCheckResourceAttr(resourceName, "attributes.0.name", "Incentro Test Shipping Method"),
					resource.TestCheckResourceAttr(resourceName, "attributes.0.scheme", "flat"),
					resource.TestCheckResourceAttr(resourceName, "attributes.0.currency_code", "EUR"),
					resource.TestCheckResourceAttr(resourceName, "attributes.0.price_amount_cents", "1000"),
					resource.TestCheckResourceAttr(resourceName, "attributes.0.free_over_amount_cents", "10000"),
					resource.TestCheckResourceAttr(resourceName, "attributes.0.min_weight", "0.5"),
					resource.TestCheckResourceAttr(resourceName, "attributes.0.max_weight", "10"),
					resource.TestCheckResourceAttr(resourceName, "attributes.0.unit_of_weight", "kg"),
					resource.TestCheckResourceAttr(resourceName, "attributes.0.metadata.foo", "bar"),
					resource.TestCheckResourceAttr(resourceName, "attributes.0.enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "attributes.0.use_subtotal", "false"),
				),
			},
			{
				Config: testAccShippingMethodUpdate(resourceName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "attributes.0.name", "Incentro Test Shipping Method Updated"),
					resource.TestCheckResourceAttr(resourceName, "attributes.0.scheme", "external"),
					resource.TestCheckResourceAttr(resourceName, "attributes.0.currency_code", "CHF"),
					resource.TestCheckResourceAttr(resourceName, "attributes.0.external_prices_url", "https://api.commercelayer.io/v1/prices"),
					resource.TestCheckResourceAttr(resourceName, "attributes.0.price_amount_cents", "1"),
					resource.TestCheckResourceAttr(resourceName, "attributes.0.free_over_amount_cents", "1"),
					resource.TestCheckResourceAttr(resourceName, "attributes.0.min_weight", "1"),
					resource.TestCheckResourceAttr(resourceName, "attributes.0.max_weight", "20"),
					resource.TestCheckResourceAttr(resourceName, "attributes.0.unit_of_weight", "oz"),
					resource.TestCheckResourceAttr(resourceName, "attributes.0.metadata.bar", "foo"),
					resource.TestCheckResourceAttr(resourceName, "attributes.0.enabled", "false"),
					resource.TestCheckResourceAttr(resourceName, "attributes.0.use_subtotal", "true"),
				),
			},
		},
	})
}

func testAccShippingMethodCreate(testName string) string {
	return hclTemplate(`
		resource "commercelayer_shipping_method" "incentro_shipping_method" {
		  attributes {
			name                   = "Incentro Test Shipping Method"
			scheme                 = "flat"
			currency_code          = "EUR"
			price_amount_cents     = 1000
			free_over_amount_cents = 10000
			min_weight             = 0.50
			max_weight             = 10
			unit_of_weight         = "kg"
			metadata               = {
			  foo : "bar"
		 	  testName: "{{.testName}}"
			}
		  }
		}
	`, map[string]any{"testName": testName})
}

func testAccShippingMethodUpdate(testName string) string {
	return hclTemplate(`
		resource "commercelayer_shipping_method" "incentro_shipping_method" {
		  attributes {
			name                   = "Incentro Test Shipping Method Updated"
			scheme                 = "external"
			currency_code          = "CHF"
			external_prices_url    = "https://api.commercelayer.io/v1/prices"
			price_amount_cents     = 1
			free_over_amount_cents = 1
			min_weight             = 1
			max_weight             = 20
			unit_of_weight         = "oz"
			enabled                = false
			use_subtotal           = true
			metadata               = {
			  bar : "foo"
		 	  testName: "{{.testName}}"
			}
		  }
		}
	`, map[string]any{"testName": testName})
}
