package commercelayer

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	commercelayer "github.com/incentro-dc/go-commercelayer-sdk/api"
	"testing"
)

func testAccCheckExternalGatewayDestroy(s *terraform.State) error {
	client := testAccProviderCommercelayer.Meta().(*commercelayer.APIClient)

	for _, rs := range s.RootModule().Resources {
		if rs.Type == "commercelayer_external_gateway" {
			_, resp, err := client.ExternalGatewaysApi.GETExternalGatewaysExternalGatewayId(context.Background(), rs.Primary.ID).Execute()
			if resp.StatusCode == 404 {
				fmt.Printf("commercelayer_external_gateway with id %s has been removed\n", rs.Primary.ID)
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

func TestAccExternalGateway_basic(t *testing.T) {
	resourceName := "commercelayer_external_gateway.incentro_external_gateway"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckExternalGatewayDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccExternalGatewayCreate(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "type", externalGatewayType),
					resource.TestCheckResourceAttr(resourceName, "attributes.0.name", "incentro_external_gateway"),
					resource.TestCheckResourceAttr(resourceName, "attributes.0.metadata.foo", "bar"),
					resource.TestCheckResourceAttr(resourceName, "attributes.0.authorize_url", "https://example.com"),
					resource.TestCheckResourceAttr(resourceName, "attributes.0.capture_url", "https://example.com"),
					resource.TestCheckResourceAttr(resourceName, "attributes.0.void_url", "https://example.com"),
					resource.TestCheckResourceAttr(resourceName, "attributes.0.refund_url", "https://example.com"),
					resource.TestCheckResourceAttr(resourceName, "attributes.0.token_url", "https://example.com"),
				),
			},
			{
				Config: testAccExternalGatewayUpdateWithUrls(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "attributes.0.name", "incentro_external_gateway_changed"),
					resource.TestCheckResourceAttr(resourceName, "attributes.0.metadata.bar", "foo"),
					resource.TestCheckResourceAttr(resourceName, "attributes.0.authorize_url", "https://foo.com"),
					resource.TestCheckResourceAttr(resourceName, "attributes.0.capture_url", "https://foo.com"),
					resource.TestCheckResourceAttr(resourceName, "attributes.0.void_url", "https://foo.com"),
					resource.TestCheckResourceAttr(resourceName, "attributes.0.refund_url", "https://foo.com"),
					resource.TestCheckResourceAttr(resourceName, "attributes.0.token_url", "https://foo.com"),
				),
			},
			{
				Config: testAccExternalGatewayUpdateWithoutUrls(),
				Check: resource.ComposeTestCheckFunc(
					//TODO: check how to do custom value checks
					resource.TestCheckResourceAttr(resourceName, "attributes.0.authorize_url", ""),
					resource.TestCheckResourceAttr(resourceName, "attributes.0.capture_url", ""),
					resource.TestCheckResourceAttr(resourceName, "attributes.0.void_url", ""),
					resource.TestCheckResourceAttr(resourceName, "attributes.0.refund_url", ""),
					resource.TestCheckResourceAttr(resourceName, "attributes.0.token_url", ""),
				),
			},
		},
	})
}

func testAccExternalGatewayCreate() string {
	return hclTemplate(`
	resource "commercelayer_external_gateway" "incentro_external_gateway" {
		  attributes {
			name          = "incentro_external_gateway"
			authorize_url = "https://example.com"
			capture_url = "https://example.com"
			void_url = "https://example.com"
			refund_url = "https://example.com"
			token_url = "https://example.com"
			metadata = {
			  foo : "bar"
			}
		  }
		}
	`, map[string]any{})
}

func testAccExternalGatewayUpdateWithUrls() string {
	return hclTemplate(`
		resource "commercelayer_external_gateway" "incentro_external_gateway" {
		  attributes {
			name          = "incentro_external_gateway_changed"
			authorize_url = "https://foo.com"
			capture_url = "https://foo.com"
			void_url = "https://foo.com"
			refund_url = "https://foo.com"
			token_url = "https://foo.com"
			metadata = {
			  bar : "foo"
			}
		  }
		}
	`, map[string]any{})
}

func testAccExternalGatewayUpdateWithoutUrls() string {
	return hclTemplate(`
		resource "commercelayer_external_gateway" "incentro_external_gateway" {
		  attributes {
			name          = "incentro_external_gateway_changed"
		  }
		}
	`, map[string]any{})
}
