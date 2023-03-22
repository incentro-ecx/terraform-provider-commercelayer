package commercelayer

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	commercelayer "github.com/incentro-dc/go-commercelayer-sdk/api"
)

func testAccCheckAdyenGatewayDestroy(s *terraform.State) error {
	client := testAccProviderCommercelayer.Meta().(*commercelayer.APIClient)

	for _, rs := range s.RootModule().Resources {
		if rs.Type == "commercelayer_adyen_gateway" {
			_, resp, err := client.AdyenGatewaysApi.
				GETAdyenGatewaysAdyenGatewayId(context.Background(), rs.Primary.ID).Execute()
			if resp.StatusCode == 404 {
				fmt.Printf("commercelayer_adyen_gateway with id %s has been removed\n", rs.Primary.ID)
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

func (s *AcceptanceSuite) TestAccAdyenGateway_basic() {
	resourceName := "commercelayer_adyen_gateway.incentro_adyen_gateway"

	resource.Test(s.T(), resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(s)
		},
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckAdyenGatewayDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAdyenGatewayCreate(resourceName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "type", adyenGatewaysType),
					resource.TestCheckResourceAttr(resourceName, "attributes.0.name", "Incentro Adyen Gateway"),
					resource.TestCheckResourceAttr(resourceName, "attributes.0.metadata.foo", "bar"),
				),
			},
			{
				Config: testAccAdyenGatewayUpdate(resourceName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "attributes.0.name", "Incentro Adyen Gateway Changed"),
					resource.TestCheckResourceAttr(resourceName, "attributes.0.metadata.bar", "foo"),
				),
			},
		},
	})
}

func testAccAdyenGatewayCreate(testName string) string {
	return hclTemplate(`
		resource "commercelayer_adyen_gateway" "incentro_adyen_gateway" {
           attributes {
			name                   = "Incentro Adyen Gateway"
			merchant_account       = "xxxx-yyyy-zzzz"
			api_key       		   = "xxxx-yyyy-zzzz"
			public_key       	   = "xxxx-yyyy-zzzz"
			live_url_prefix        = "1797a841fbb37ca7-AdyenDemo"

			metadata = {
				foo: "bar"
				testName: "{{.testName}}"
    		}
  		}
	}
`, map[string]any{"testName": testName})
}

func testAccAdyenGatewayUpdate(testName string) string {
	return hclTemplate(`
		resource "commercelayer_adyen_gateway" "incentro_adyen_gateway" {
           attributes {
			name                   = "Incentro Adyen Gateway Changed"
			merchant_account       = "xxxx-yyyy-zzzz"
			api_key       		   = "xxxx-yyyy-zzzz"
			public_key       	   = "xxxx-yyyy-zzzz"
			live_url_prefix        = "1797a841fbb37ca7-AdyenDemo"

			metadata = {
				bar: "foo"
				testName: "{{.testName}}"
    		}
  		}
	}
`, map[string]any{"testName": testName})
}
