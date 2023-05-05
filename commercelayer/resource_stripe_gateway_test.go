package commercelayer

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	commercelayer "github.com/incentro-dc/go-commercelayer-sdk/api"
)

func testAccCheckStripeGatewayDestroy(s *terraform.State) error {
	client := testAccProviderCommercelayer.Meta().(*commercelayer.APIClient)

	for _, rs := range s.RootModule().Resources {
		if rs.Type == "commercelayer_stripe_gateway" {
			_, resp, err := client.StripeGatewaysApi.
				GETStripeGatewaysStripeGatewayId(context.Background(), rs.Primary.ID).Execute()
			if resp.StatusCode == 404 {
				fmt.Printf("commercelayer_stripe_gateway with id %s has been removed\n", rs.Primary.ID)
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
					resource.TestCheckResourceAttr(resourceName, "type", stripeGatewaysType),
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
			name        	= "Incentro Stripe Gateway"
			login       	= "xxxx-yyyy-zzzz"
			publishable_key = "aaaa-bbbb-cccc"

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
			name        	= "Incentro Stripe Gateway Changed"
			login       	= "xxxx-yyyy-zzzz"
			publishable_key = "aaaa-bbbb-cccc"

			metadata = {
				bar: "foo"
				testName: "{{.testName}}"
    		}
  		}
	}
`, map[string]any{"testName": testName})
}
