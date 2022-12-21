package commercelayer

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	commercelayer "github.com/incentro-dc/go-commercelayer-sdk/api"
	"net/http"
	"strings"
)

func testAccCheckStockLocationDestroy(s *terraform.State) error {
	client := testAccProviderCommercelayer.Meta().(*commercelayer.APIClient)

	for _, rs := range s.RootModule().Resources {
		if rs.Type == "commercelayer_stock_location" {
			err := retryRemoval(10, func() (*http.Response, error) {
				_, resp, err := client.StockLocationsApi.GETStockLocationsStockLocationId(context.Background(), rs.Primary.ID).
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

func (s *AcceptanceSuite) TestAccStockLocation_basic() {
	resourceName := "commercelayer_stock_location.incentro_stock_location"

	resource.Test(s.T(), resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(s)
		},
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckStockLocationDestroy,
		Steps: []resource.TestStep{
			{
				Config: strings.Join([]string{testAccAddressCreate(resourceName), testAccStockLocationCreate(resourceName)}, "\n"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "type", stockLocationType),
					resource.TestCheckResourceAttr(resourceName, "attributes.0.name", "Incentro Warehouse Stock Location"),
					resource.TestCheckResourceAttr(resourceName, "attributes.0.label_format", "PNG"),
					resource.TestCheckResourceAttr(resourceName, "attributes.0.suppress_etd", "true"),
					resource.TestCheckResourceAttr(resourceName, "attributes.0.metadata.foo", "bar"),
				),
			},
			{
				Config: strings.Join([]string{testAccAddressCreate(resourceName), testAccStockLocationUpdate(resourceName)}, "\n"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "attributes.0.name", "Incentro Warehouse Stock Location Updated"),
					resource.TestCheckResourceAttr(resourceName, "attributes.0.label_format", "PDF"),
					resource.TestCheckResourceAttr(resourceName, "attributes.0.suppress_etd", "false"),
					resource.TestCheckResourceAttr(resourceName, "attributes.0.metadata.bar", "foo"),
				),
			},
		},
	})
}

func testAccStockLocationCreate(testName string) string {
	return hclTemplate(`
		resource "commercelayer_stock_location" "incentro_stock_location" {
		  attributes {
			name         = "Incentro Warehouse Stock Location"
			label_format = "PNG"
			suppress_etd = true
			metadata     = {
			  foo : "bar"
		 	  testName: "{{.testName}}"
			}
		  }
		
		  relationships {
			address_id = commercelayer_address.incentro_address.id
		  }
		}
	`, map[string]any{"testName": testName})
}

func testAccStockLocationUpdate(testName string) string {
	return hclTemplate(`
		resource "commercelayer_stock_location" "incentro_stock_location" {
		  attributes {
			name         = "Incentro Warehouse Stock Location Updated"
			label_format = "PDF"
			suppress_etd = false
			metadata     = {
			  bar : "foo"
		 	  testName: "{{.testName}}"
			}
		  }
		
		  relationships {
			address_id = commercelayer_address.incentro_address.id
		  }
		}
	`, map[string]any{"testName": testName})
}
