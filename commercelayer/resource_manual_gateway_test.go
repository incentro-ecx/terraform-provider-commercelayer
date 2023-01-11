package commercelayer

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	commercelayer "github.com/incentro-dc/go-commercelayer-sdk/api"
	"net/http"
)

func testAccCheckManualGatewayDestroy(s *terraform.State) error {
	client := testAccProviderCommercelayer.Meta().(*commercelayer.APIClient)

	for _, rs := range s.RootModule().Resources {
		if rs.Type == "commercelayer_manual_gateway" {
			err := retryRemoval(10, func() (*http.Response, error) {
				_, resp, err := client.ManualGatewaysApi.
					GETManualGatewaysManualGatewayId(context.Background(), rs.Primary.ID).
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

func (s *AcceptanceSuite) TestAccManualGateway_basic() {
	resourceName := "commercelayer_manual_gateway.incentro_manual_gateway"

	resource.Test(s.T(), resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(s)
		},
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckManualGatewayDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccManualGatewayCreate(resourceName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "type", manualGatewaysType),
					resource.TestCheckResourceAttr(resourceName, "attributes.0.name", "Incentro Manual Gateway"),
					resource.TestCheckResourceAttr(resourceName, "attributes.0.metadata.foo", "bar"),
				),
			},
			{
				Config: testAccManualGatewayUpdate(resourceName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "attributes.0.name", "Incentro Manual Gateway Changed"),
					resource.TestCheckResourceAttr(resourceName, "attributes.0.metadata.bar", "foo"),
				),
			},
		},
	})
}

func testAccManualGatewayCreate(testName string) string {
	return hclTemplate(`
		resource "commercelayer_manual_gateway" "incentro_manual_gateway" {
           attributes {
			name                   = "Incentro Manual Gateway"
			metadata = {
				foo: "bar"
				testName: "{{.testName}}"
    		}
  		}
	}
`, map[string]any{"testName": testName})
}

func testAccManualGatewayUpdate(testName string) string {
	return hclTemplate(`
		resource "commercelayer_manual_gateway" "incentro_manual_gateway" {
           attributes {
			name                   = "Incentro Manual Gateway Changed"
			metadata = {
				bar: "foo"
				testName: "{{.testName}}"
    		}
  		}
	}
`, map[string]any{"testName": testName})
}
