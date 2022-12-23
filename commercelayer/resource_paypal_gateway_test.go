package commercelayer

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	commercelayer "github.com/incentro-dc/go-commercelayer-sdk/api"
	"net/http"
)

func testAccCheckPaypalGatewayDestroy(s *terraform.State) error {
	client := testAccProviderCommercelayer.Meta().(*commercelayer.APIClient)

	for _, rs := range s.RootModule().Resources {
		if rs.Type == "commercelayer_paypal_gateway" {
			err := retryRemoval(10, func() (*http.Response, error) {
				_, resp, err := client.PaypalGatewaysApi.
					GETPaypalGatewaysPaypalGatewayId(context.Background(), rs.Primary.ID).
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
