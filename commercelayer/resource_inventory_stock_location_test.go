package commercelayer

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	commercelayer "github.com/incentro-dc/go-commercelayer-sdk/api"
	"net/http"
	"strings"
)

func testAccCheckInventoryStockLocationDestroy(s *terraform.State) error {
	client := testAccProviderCommercelayer.Meta().(*commercelayer.APIClient)

	for _, rs := range s.RootModule().Resources {
		if rs.Type == "commercelayer_inventory_stock_location" {
			err := retryRemoval(10, func() (*http.Response, error) {
				_, resp, err := client.InventoryStockLocationsApi.
					GETInventoryStockLocationsInventoryStockLocationId(context.Background(), rs.Primary.ID).
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

func (s *AcceptanceSuite) TestAccInventoryStockLocation_basic() {
	resourceName := "commercelayer_inventory_stock_location.incentro_inventory_stock_location"

	resource.Test(s.T(), resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(s)
		},
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckInventoryStockLocationDestroy,
		Steps: []resource.TestStep{
			{
				Config: strings.Join([]string{
					testAccAddressCreate(resourceName),
					testAccInventoryModelCreate(resourceName),
					testAccStockLocationCreate(resourceName),
					testAccInventoryStockLocationCreate(resourceName),
				}, "\n"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "attributes.0.priority", "1"),
					resource.TestCheckResourceAttr(resourceName, "attributes.0.on_hold", "true"),
				),
			},
			{
				Config: strings.Join([]string{
					testAccAddressCreate(resourceName),
					testAccInventoryModelCreate(resourceName),
					testAccStockLocationCreate(resourceName),
					testAccInventoryStockLocationUpdate(resourceName),
				}, "\n"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "attributes.0.priority", "2"),
					resource.TestCheckResourceAttr(resourceName, "attributes.0.on_hold", "false"),
				),
			},
		},
	})
}

func testAccInventoryStockLocationCreate(testName string) string {
	return hclTemplate(`
		resource "commercelayer_inventory_stock_location" "incentro_inventory_stock_location" {
		  attributes {
			priority = 1
			on_hold  = true
			metadata = {
			  testName: "{{.testName}}"
			}
		  }

		  relationships {
			inventory_model_id = commercelayer_inventory_model.incentro_inventory_model.id
			stock_location_id  = commercelayer_stock_location.incentro_stock_location.id
		  }
		}
	`, map[string]any{"testName": testName})
}

func testAccInventoryStockLocationUpdate(testName string) string {
	return hclTemplate(`
		resource "commercelayer_inventory_stock_location" "incentro_inventory_stock_location" {
		  attributes {
			priority = 2
			on_hold  = false
			metadata = {
			  testName: "{{.testName}}"
			}
		  }

		  relationships {
			inventory_model_id = commercelayer_inventory_model.incentro_inventory_model.id
			stock_location_id  = commercelayer_stock_location.incentro_stock_location.id
		  }
		}
	`, map[string]any{"testName": testName})
}
