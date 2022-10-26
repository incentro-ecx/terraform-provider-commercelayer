package commercelayer

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	commercelayer "github.com/incentro-dc/go-commercelayer-sdk/api"
	"net/http"
	"testing"
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

func TestAccInventoryModel_basic(t *testing.T) {
	resourceName := "commercelayer_inventory_model.incentro_inventory_model"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckInventoryModelDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccInventoryModelCreate(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "attributes.0.name", "Incentro Inventory Model"),
					resource.TestCheckResourceAttr(resourceName, "attributes.0.stock_locations_cutoff", "1"),
					resource.TestCheckResourceAttr(resourceName, "attributes.0.strategy", "no_split"),
				),
			},
			{
				Config: testAccInventoryModelUpdate(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "attributes.0.name", "Incentro Inventory Model Changed"),
					resource.TestCheckResourceAttr(resourceName, "attributes.0.stock_locations_cutoff", "2"),
					resource.TestCheckResourceAttr(resourceName, "attributes.0.strategy", "split_shipments"),
				),
			},
		},
	})
}

func testAccInventoryModelCreate() string {
	return hclTemplate(`
		resource "commercelayer_inventory_model" "incentro_inventory_model" {
		  attributes {
			name                   = "Incentro Inventory Model"
			stock_locations_cutoff = 1
			strategy               = "no_split"
		  }
		}
	`, map[string]any{})
}

func testAccInventoryModelUpdate() string {
	return hclTemplate(`
		resource "commercelayer_inventory_model" "incentro_inventory_model" {
		  attributes {
			name                   = "Incentro Inventory Model Changed"
			stock_locations_cutoff = 2
			strategy               = "split_shipments"
		  }
		}
	`, map[string]any{})
}
