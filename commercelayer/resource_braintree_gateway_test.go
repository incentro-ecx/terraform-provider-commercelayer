package commercelayer

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	commercelayer "github.com/incentro-dc/go-commercelayer-sdk/api"
)

func testAccCheckBraintreeGatewayDestroy(s *terraform.State) error {
	client := testAccProviderCommercelayer.Meta().(*commercelayer.APIClient)

	for _, rs := range s.RootModule().Resources {
		if rs.Type == "commercelayer_braintree_gateway" {
			_, resp, err := client.BraintreeGatewaysApi.
				GETBraintreeGatewaysBraintreeGatewayId(context.Background(), rs.Primary.ID).Execute()
			if resp.StatusCode == 404 {
				fmt.Printf("commercelayer_braintree_gateway with id %s has been removed\n", rs.Primary.ID)
				continue
			}
			if err != nil {
				return err
			}

			return fmt.Errorf("received response code with status %d", resp.StatusCode)
		}

	}
	return nil
}

func (s *AcceptanceSuite) TestAccBraintreeGateway_basic() {
	resourceName := "commercelayer_braintree_gateway.incentro_braintree_gateway"

	resource.Test(s.T(), resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(s)
		},
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckBraintreeGatewayDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccBraintreeGatewayCreate(resourceName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "type", braintreeGatewaysType),
					resource.TestCheckResourceAttr(resourceName, "attributes.0.name", "Incentro Braintree Gateway"),
					resource.TestCheckResourceAttr(resourceName, "attributes.0.metadata.foo", "bar"),
				),
			},
			{
				Config: testAccBraintreeGatewayUpdate(resourceName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "attributes.0.name", "Incentro Braintree Gateway Changed"),
					resource.TestCheckResourceAttr(resourceName, "attributes.0.metadata.bar", "foo"),
				),
			},
		},
	})
}

func testAccBraintreeGatewayCreate(testName string) string {
	return hclTemplate(`
		resource "commercelayer_braintree_gateway" "incentro_braintree_gateway" {
           attributes {
			name                   = "Incentro Braintree Gateway"
			merchant_account_id    = "xxxx-yyyy-zzzz"
			merchant_id            = "xxxx-yyyy-zzzz"
			public_key             = "xxxx-yyyy-zzzz"
			private_key            = "xxxx-yyyy-zzzz"

			metadata = {
				foo: "bar"
				testName: "{{.testName}}"
    		}
  		}
	}
`, map[string]any{"testName": testName})
}

func testAccBraintreeGatewayUpdate(testName string) string {
	return hclTemplate(`
		resource "commercelayer_braintree_gateway" "incentro_braintree_gateway" {
           attributes {
			name                        = "Incentro Braintree Gateway Changed"
			merchant_account_id         = "xxxx-yyyy-zzzz"
			merchant_id                 = "xxxx-yyyy-zzzz"
			public_key                  = "xxxx-yyyy-zzzz"
			private_key                 = "xxxx-yyyy-zzzz"

			metadata = {
				bar: "foo"
				testName: "{{.testName}}"
    		}
  		}
	}
`, map[string]any{"testName": testName})
}
