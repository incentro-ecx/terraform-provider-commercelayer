package commercelayer

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	commercelayer "github.com/incentro-dc/go-commercelayer-sdk/api"
	"net/http"
)

func testAccCheckPaymentMethodDestroy(s *terraform.State) error {
	client := testAccProviderCommercelayer.Meta().(*commercelayer.APIClient)

	for _, rs := range s.RootModule().Resources {
		if rs.Type == "commercelayer_payment_method" {
			err := retryRemoval(10, func() (*http.Response, error) {
				_, resp, err := client.PaymentMethodsApi.GETPaymentMethodsPaymentMethodId(context.Background(), rs.Primary.ID).
					Execute()
				return resp, err
			})
			if err != nil {
				return err
			}
		}

	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type == "commercelayer_market" {
			err := retryRemoval(10, func() (*http.Response, error) {
				_, resp, err := client.MarketsApi.
					GETMarketsMarketId(context.Background(), rs.Primary.ID).
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

func (s *AcceptanceSuite) TestAccPaymentMethod_basic() {
	resourceName := "commercelayer_payment_method.incentro_payment_method"

	resource.Test(s.T(), resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(s)
		},
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckPaymentMethodDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccPaymentMethodCreate(resourceName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "type", paymentMethodType),
				),
			},
			{
				Config: testAccPaymentMethodUpdate(resourceName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "attributes.0.name", "Incentro Payment Method Updated"),
				),
			},
		},
	})
}

// TODO: add payment_gateway_id to Template body
func testAccPaymentMethodCreate(testName string) string {
	return hclTemplate(`
		resource "commercelayer_payment_method" "incentro_payment_method" {
		  attributes {
      		payment_source_type   = "CreditCard",
			currency_code          = "EUR"
			price_amount_cents     = 1000
			metadata               = {
			  foo : "bar"
		 	  testName: "{{.testName}}"
			}
		  }

		  relationships {
			market_id = commercelayer_market.incentro_market.id
		  }
		}
	`, map[string]any{"testName": testName})
}

// TODO: add payment_gateway_id to Template body
func testAccPaymentMethodUpdate(testName string) string {
	return hclTemplate(`
		resource "commercelayer_payment_method" "incentro_payment_method" {
		  attributes {
      		payment_source_type    = "CreditCard"
			currency_code          = "EUR"
			price_amount_cents     = 1000
			metadata               = {
			  foo : "bar"
		 	  testName: "{{.testName}}"
			}
		  }
  		  relationships {
			market_id = commercelayer_market.incentro_market.id
		  }
		}
	`, map[string]any{"testName": testName})
}
