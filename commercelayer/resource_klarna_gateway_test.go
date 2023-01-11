package commercelayer

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	commercelayer "github.com/incentro-dc/go-commercelayer-sdk/api"
)

func testAccCheckKlarnaGatewayDestroy(s *terraform.State) error {
	client := testAccProviderCommercelayer.Meta().(*commercelayer.APIClient)

	for _, rs := range s.RootModule().Resources {
		if rs.Type == "commercelayer_klarna_gateway" {
			_, resp, err := client.KlarnaGatewaysApi.
				GETKlarnaGatewaysKlarnaGatewayId(context.Background(), rs.Primary.ID).Execute()
			if resp.StatusCode == 404 {
				fmt.Printf("commercelayer_klarna_gateway with id %s has been removed\n", rs.Primary.ID)
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

func (s *AcceptanceSuite) TestAccKlarnaGateway_basic() {
	resourceName := "commercelayer_klarna_gateway.incentro_klarna_gateway"

	resource.Test(s.T(), resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(s)
		},
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckKlarnaGatewayDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccKlarnaGatewayCreate(resourceName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "type", klarnaGatewaysType),
					resource.TestCheckResourceAttr(resourceName, "attributes.0.name", "Incentro Klarna Gateway"),
					resource.TestCheckResourceAttr(resourceName, "attributes.0.metadata.foo", "bar"),
				),
			},
			{
				Config: testAccKlarnaGatewayUpdate(resourceName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "attributes.0.name", "Incentro Klarna Gateway Changed"),
					resource.TestCheckResourceAttr(resourceName, "attributes.0.metadata.bar", "foo"),
				),
			},
		},
	})
}

func testAccKlarnaGatewayCreate(testName string) string {
	return hclTemplate(`
		resource "commercelayer_klarna_gateway" "incentro_klarna_gateway" {
           attributes {
			name                   = "Incentro Klarna Gateway"
			country_code              = "EU"
			api_key              = "xxxx-yyyy-zzzz"
			api_secret          = "xxxx-yyyy-zzzz"

			metadata = {
				foo: "bar"
				testName: "{{.testName}}"
    		}
  		}
	}
`, map[string]any{"testName": testName})
}

func testAccKlarnaGatewayUpdate(testName string) string {
	return hclTemplate(`
		resource "commercelayer_klarna_gateway" "incentro_klarna_gateway" {
           attributes {
			name                   = "Incentro Klarna Gateway Changed"
			country_code              = "EU"
			api_key              = "xxxx-yyyy-zzzz"
			api_secret          = "xxxx-yyyy-zzzz"

			metadata = {
				bar: "foo"
				testName: "{{.testName}}"
    		}
  		}
	}
`, map[string]any{"testName": testName})
}
