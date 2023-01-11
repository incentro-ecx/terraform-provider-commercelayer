package commercelayer

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	commercelayer "github.com/incentro-dc/go-commercelayer-sdk/api"
)

func testAccCheckPaypalGatewayDestroy(s *terraform.State) error {
	client := testAccProviderCommercelayer.Meta().(*commercelayer.APIClient)

	for _, rs := range s.RootModule().Resources {
		if rs.Type == "commercelayer_paypal_gateway" {
			_, resp, err := client.PaypalGatewaysApi.
				GETPaypalGatewaysPaypalGatewayId(context.Background(), rs.Primary.ID).Execute()
			if resp.StatusCode == 404 {
				fmt.Printf("commercelayer_paypal_gateway with id %s has been removed\n", rs.Primary.ID)
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

func (s *AcceptanceSuite) TestAccPaypalGateway_basic() {
	resourceName := "commercelayer_paypal_gateway.incentro_paypal_gateway"

	resource.Test(s.T(), resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(s)
		},
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckPaypalGatewayDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccPaypalGatewayCreate(resourceName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "type", paypalGatewaysType),
					resource.TestCheckResourceAttr(resourceName, "attributes.0.name", "Incentro Paypal Gateway"),
					resource.TestCheckResourceAttr(resourceName, "attributes.0.metadata.foo", "bar"),
				),
			},
			{
				Config: testAccPaypalGatewayUpdate(resourceName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "attributes.0.name", "Incentro Paypal Gateway Changed"),
					resource.TestCheckResourceAttr(resourceName, "attributes.0.metadata.bar", "foo"),
				),
			},
		},
	})
}

func testAccPaypalGatewayCreate(testName string) string {
	return hclTemplate(`
		resource "commercelayer_paypal_gateway" "incentro_paypal_gateway" {
           attributes {
			name                   = "Incentro Paypal Gateway"
			client_id              = "xxxx-yyyy-zzzz"
			client_secret          = "xxxx-yyyy-zzzz"

			metadata = {
				foo: "bar"
				testName: "{{.testName}}"
    		}
  		}
	}
`, map[string]any{"testName": testName})
}

func testAccPaypalGatewayUpdate(testName string) string {
	return hclTemplate(`
		resource "commercelayer_paypal_gateway" "incentro_paypal_gateway" {
           attributes {
			name                   = "Incentro Paypal Gateway Changed"
			client_id              = "xxxx-yyyy-zzzz"
			client_secret          = "xxxx-yyyy-zzzz"

			metadata = {
				bar: "foo"
				testName: "{{.testName}}"
    		}
  		}
	}
`, map[string]any{"testName": testName})
}
