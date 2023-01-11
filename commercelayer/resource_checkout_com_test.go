package commercelayer

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	commercelayer "github.com/incentro-dc/go-commercelayer-sdk/api"
)

func testAccCheckCheckoutComGatewayDestroy(s *terraform.State) error {
	client := testAccProviderCommercelayer.Meta().(*commercelayer.APIClient)

	for _, rs := range s.RootModule().Resources {
		if rs.Type == "commercelayer_checkout_com_gateway" {
			_, resp, err := client.CheckoutComGatewaysApi.
				GETCheckoutComGatewaysCheckoutComGatewayId(context.Background(), rs.Primary.ID).Execute()
			if resp.StatusCode == 404 {
				fmt.Printf("commercelayer_checkout_com_gateway with id %s has been removed\n", rs.Primary.ID)
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

func (s *AcceptanceSuite) TestAccCheckoutComGateway_basic() {
	resourceName := "commercelayer_checkout_com_gateway.incentro_checkout_com_gateway"

	resource.Test(s.T(), resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(s)
		},
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckCheckoutComGatewayDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckoutComGatewayCreate(resourceName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "type", checkoutComGatewaysType),
					resource.TestCheckResourceAttr(resourceName, "attributes.0.name", "Incentro CheckoutCom Gateway"),
					resource.TestCheckResourceAttr(resourceName, "attributes.0.metadata.foo", "bar"),
				),
			},
			{
				Config: testAccCheckoutComGatewayUpdate(resourceName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "attributes.0.name", "Incentro CheckoutCom Gateway Changed"),
					resource.TestCheckResourceAttr(resourceName, "attributes.0.metadata.bar", "foo"),
				),
			},
		},
	})
}

func testAccCheckoutComGatewayCreate(testName string) string {
	return hclTemplate(`
		resource "commercelayer_checkout_com_gateway" "incentro_checkout_com_gateway" {
           attributes {
			name                   = "Incentro CheckoutCom Gateway"
			secret_key 			   = "sk_test_xxxx-yyyy-zzzz"
			public_key 			   = "pk_test_xxxx-yyyy-zzzz"

			metadata = {
				foo: "bar"
				testName: "{{.testName}}"
    		}
  		}
	}
`, map[string]any{"testName": testName})
}

func testAccCheckoutComGatewayUpdate(testName string) string {
	return hclTemplate(`
		resource "commercelayer_checkout_com_gateway" "incentro_checkout_com_gateway" {
           attributes {
			name                   = "Incentro CheckoutCom Gateway Changed"
			secret_key 			   = "sk_test_xxxx-yyyy-zzzz"
			public_key 			   = "pk_test_xxxx-yyyy-zzzz"

			metadata = {
				bar: "foo"
				testName: "{{.testName}}"
    		}
  		}
	}
`, map[string]any{"testName": testName})
}
