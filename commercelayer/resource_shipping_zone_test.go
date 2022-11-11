package commercelayer

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	commercelayer "github.com/incentro-dc/go-commercelayer-sdk/api"
	"net/http"
)

func testAccCheckShippingZoneDestroy(s *terraform.State) error {
	client := testAccProviderCommercelayer.Meta().(*commercelayer.APIClient)

	for _, rs := range s.RootModule().Resources {
		if rs.Type == "commercelayer_shipping_zone" {
			err := retryRemoval(10, func() (*http.Response, error) {
				_, resp, err := client.ShippingZonesApi.GETShippingZonesShippingZoneId(context.Background(), rs.Primary.ID).
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

func (s *AcceptanceSuite) TestAccShippingZone_basic() {
	resourceName := "commercelayer_shipping_zone.incentro_shipping_zone"

	resource.Test(s.T(), resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(s)
		},
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckShippingZoneDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccShippingZoneCreate(resourceName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "type", shippingZoneType),
					resource.TestCheckResourceAttr(resourceName, "attributes.0.name", "Incentro Shipping Zone"),
					resource.TestCheckResourceAttr(resourceName, "attributes.0.country_code_regex", ".*"),
					resource.TestCheckResourceAttr(resourceName, "attributes.0.not_country_code_regex", "[^i*&2@]"),
					resource.TestCheckResourceAttr(resourceName, "attributes.0.state_code_regex", "^dog"),
					resource.TestCheckResourceAttr(resourceName, "attributes.0.not_state_code_regex", "//[^\r\n]*[\r\n]"),
					resource.TestCheckResourceAttr(resourceName, "attributes.0.zip_code_regex", "[a-zA-Z]{2,4}"),
					resource.TestCheckResourceAttr(resourceName, "attributes.0.not_zip_code_regex", ".+"),
					resource.TestCheckResourceAttr(resourceName, "attributes.0.metadata.foo", "bar"),
				),
			},
			{
				Config: testAccShippingZoneUpdate(resourceName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "attributes.0.name", "Incentro Shipping Zone Updated"),
					resource.TestCheckResourceAttr(resourceName, "attributes.0.country_code_regex", ".+"),
					resource.TestCheckResourceAttr(resourceName, "attributes.0.not_country_code_regex", "[^i*&2@]G"),
					resource.TestCheckResourceAttr(resourceName, "attributes.0.state_code_regex", "^cat"),
					resource.TestCheckResourceAttr(resourceName, "attributes.0.not_state_code_regex", "//[^\r\n]"),
					resource.TestCheckResourceAttr(resourceName, "attributes.0.zip_code_regex", "[a-z]{1,2}"),
					resource.TestCheckResourceAttr(resourceName, "attributes.0.not_zip_code_regex", ".*"),
					resource.TestCheckResourceAttr(resourceName, "attributes.0.metadata.bar", "foo"),
				),
			},
		},
	})
}

func testAccShippingZoneCreate(testName string) string {
	return hclTemplate(`
		resource "commercelayer_shipping_zone" "incentro_shipping_zone" {
		  attributes {
			name                   = "Incentro Shipping Zone"
			country_code_regex     = ".*"
			not_country_code_regex = "[^i*&2@]"
			state_code_regex       = "^dog"
			not_state_code_regex   = "//[^\r\n]*[\r\n]"
			zip_code_regex         = "[a-zA-Z]{2,4}"
			not_zip_code_regex     = ".+"
			metadata               = {
			  foo : "bar"
		 	  testName: "{{.testName}}"
			}
		  }
		}
	`, map[string]any{"testName": testName})
}

func testAccShippingZoneUpdate(testName string) string {
	return hclTemplate(`
		resource "commercelayer_shipping_zone" "incentro_shipping_zone" {
		  attributes {
			name                   = "Incentro Shipping Zone Updated"
			country_code_regex     = ".+"
			not_country_code_regex = "[^i*&2@]G"
			state_code_regex       = "^cat"
			not_state_code_regex   = "//[^\r\n]"
			zip_code_regex         = "[a-z]{1,2}"
			not_zip_code_regex     = ".*"
			metadata               = {
			  bar : "foo"
		 	  testName: "{{.testName}}"
			}
		  }
		}
	`, map[string]any{"testName": testName})
}
