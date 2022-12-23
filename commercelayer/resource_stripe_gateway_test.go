package commercelayer

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	commercelayer "github.com/incentro-dc/go-commercelayer-sdk/api"
	"net/http"
)

func testAccCheckStripeGatewayDestroy(s *terraform.State) error {
	client := testAccProviderCommercelayer.Meta().(*commercelayer.APIClient)

	for _, rs := range s.RootModule().Resources {
		if rs.Type == "commercelayer_stripe_gateway" {
			err := retryRemoval(10, func() (*http.Response, error) {
				_, resp, err := client.StripeGatewaysApi.
					GETStripeGatewaysStripeGatewayId(context.Background(), rs.Primary.ID).
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

func (s *AcceptanceSuite) TestAccStripeGateway_basic() {
	resourceName := "commercelayer_stripe_gateway.incentro_stripe_gateway"

	resource.Test(s.T(), resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(s)
		},
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckStripeGatewayDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccStripeGatewayCreate(resourceName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "attributes.0.name", "Incentro Stripe Gateway"),
					resource.TestCheckResourceAttr(resourceName, "attributes.0.metadata.foo", "bar"),
				),
			},
			{
				Config: testAccStripeGatewayUpdate(resourceName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "attributes.0.name", "Incentro Stripe Gateway Changed"),
					resource.TestCheckResourceAttr(resourceName, "attributes.0.metadata.bar", "foo"),
				),
			},
		},
	})
}

func testAccStripeGatewayCreate(testName string) string {
	return hclTemplate(`
		resource "commercelayer_stripe_gateway" "incentro_stripe_gateway" {
           attributes {
			name                   = "Incentro Stripe Gateway"
			login                  = "SecretPassword"
			metadata = {
				foo: "bar"
				testName: "{{.testName}}"
    		}
  		}
	}
`, map[string]any{"testName": testName})
}

func testAccStripeGatewayUpdate(testName string) string {
	return hclTemplate(`
		resource "commercelayer_stripe_gateway" "incentro_stripe_gateway" {
           attributes {
			name                   = "Incentro Stripe Gateway Changed"
			metadata = {
				bar: "foo"
				testName: "{{.testName}}"
    		}
  		}
	}
`, map[string]any{"testName": testName})
}
