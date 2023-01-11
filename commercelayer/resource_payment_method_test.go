package commercelayer

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	commercelayer "github.com/incentro-dc/go-commercelayer-sdk/api"
	"net/http"
	"strings"
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
		if rs.Type == "commercelayer_payment_method" {
			err := retryRemoval(10, func() (*http.Response, error) {
				_, resp, err := client.PaymentMethodsApi.
					GETPaymentMethodsPaymentMethodId(context.Background(), rs.Primary.ID).
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
				Config: strings.Join([]string{testAccPaymentMethodCreate(resourceName), testAccAdyenGatewayCreate(resourceName), testAccMarketCreate(resourceName)}, "\n"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "type", paymentMethodType),
					resource.TestCheckResourceAttr(resourceName, "attributes.0.metadata.foo", "bar"),
					resource.TestCheckResourceAttr(resourceName, "attributes.0.currency_code", "EUR"),
					resource.TestCheckResourceAttr(resourceName, "attributes.0.payment_source_type", "CreditCard"),
					resource.TestCheckResourceAttr(resourceName, "attributes.0.price_amount_cents", "0"),
				),
			},
			{
				Config: strings.Join([]string{testAccPaymentMethodCreate(resourceName), testAccAdyenGatewayCreate(resourceName), testAccMarketCreate(resourceName)}, "\n"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "attributes.0.metadata.bar", "foo"),
					resource.TestCheckResourceAttr(resourceName, "attributes.0.currency_code", "EUR"),
					resource.TestCheckResourceAttr(resourceName, "attributes.0.payment_source_type", "CreditCard"),
					resource.TestCheckResourceAttr(resourceName, "attributes.0.price_amount_cents", "0"),
				),
			},
		},
	})
}

func testAccPaymentMethodCreate(testName string) string {
	return hclTemplate(`
		resource "commercelayer_payment_method" "incentro_payment_method" {
		  attributes {
      		payment_source_type   = "CreditCard"
			currency_code          = "EUR"
			price_amount_cents     = 0
			metadata               = {
			  foo : "bar"
		 	  testName: "{{.testName}}"
			}
		  }

		  relationships {
			payment_gateway_id = commercelayer_payment_gateway.incentro_payment_gateway.id
			market_id = commercelayer_market.incentro_market.id
		  }
		}
	`, map[string]any{"testName": testName})
}

func testAccPaymentMethodUpdate(testName string) string {
	return hclTemplate(`
		resource "commercelayer_payment_method" "incentro_payment_method" {
		  attributes {
      		payment_source_type    = "CreditCard"
			currency_code          = "EUR"
			price_amount_cents     = 0
			metadata               = {
			  bar : "foo"
		 	  testName: "{{.testName}}"
			}
		  }
  		  relationships {
			payment_gateway_id = commercelayer_payment_gateway.incentro_payment_gateway.id
			market_id = commercelayer_market.incentro_market.id
		  }
		}
	`, map[string]any{"testName": testName})
}
