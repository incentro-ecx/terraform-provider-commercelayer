package commercelayer

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	commercelayer "github.com/incentro-dc/go-commercelayer-sdk/api"
	"net/http"
)

func testAccCheckInventoryModelDestroy(s *terraform.State) error {
	client := testAccProviderCommercelayer.Meta().(*commercelayer.APIClient)

	for _, rs := range s.RootModule().Resources {
		if rs.Type == "commercelayer_inventory_model" {
			err := retryRemoval(10, func() (*http.Response, error) {
				_, resp, err := client.InventoryModelsApi.
					GETInventoryModelsInventoryModelId(context.Background(), rs.Primary.ID).
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

func (s *AcceptanceSuite) TestAccInventoryModel_basic() {
	resourceName := "commercelayer_inventory_model.incentro_inventory_model"

	resource.Test(s.T(), resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(s)
		},
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckInventoryModelDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccInventoryModelCreate(resourceName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "attributes.0.name", "Incentro Inventory Model"),
					resource.TestCheckResourceAttr(resourceName, "attributes.0.stock_locations_cutoff", "1"),
					resource.TestCheckResourceAttr(resourceName, "attributes.0.strategy", "no_split"),
					resource.TestCheckResourceAttr(resourceName, "attributes.0.manual_stock_decrement", "true"),
					resource.TestCheckResourceAttr(resourceName, "attributes.0.stock_reservation_cutoff", "4000"),
					resource.TestCheckResourceAttr(resourceName, "attributes.0.put_stock_transfers_on_hold", "true"),
				),
			},
			{
				Config: testAccInventoryModelUpdate(resourceName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "attributes.0.name", "Incentro Inventory Model Changed"),
					resource.TestCheckResourceAttr(resourceName, "attributes.0.stock_locations_cutoff", "2"),
					resource.TestCheckResourceAttr(resourceName, "attributes.0.strategy", "split_shipments"),
					resource.TestCheckResourceAttr(resourceName, "attributes.0.manual_stock_decrement", "false"),
					resource.TestCheckResourceAttr(resourceName, "attributes.0.stock_reservation_cutoff", "3600"),
					resource.TestCheckResourceAttr(resourceName, "attributes.0.put_stock_transfers_on_hold", "false"),
				),
			},
		},
	})
}

func testAccInventoryModelCreate(testName string) string {
	return hclTemplate(`
		resource "commercelayer_inventory_model" "incentro_inventory_model" {
		  attributes {
			name                   = "Incentro Inventory Model"
			stock_locations_cutoff = 1
			strategy               = "no_split"
			manual_stock_decrement = true
			stock_reservation_cutoff = 4000
			put_stock_transfers_on_hold = true
			metadata = {
			  testName: "{{.testName}}"
			}
		  }
		}
	`, map[string]any{"testName": testName})
}

func testAccInventoryModelUpdate(testName string) string {
	return hclTemplate(`
		resource "commercelayer_inventory_model" "incentro_inventory_model" {
		  attributes {
			name                   = "Incentro Inventory Model Changed"
			stock_locations_cutoff = 2
			strategy               = "split_shipments"
			manual_stock_decrement = false
			stock_reservation_cutoff = 3600
			put_stock_transfers_on_hold = false
			metadata = {
			  testName: "{{.testName}}"
			}
		  }
		}
	`, map[string]any{"testName": testName})
}
