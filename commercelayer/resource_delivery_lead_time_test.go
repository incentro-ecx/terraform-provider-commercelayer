package commercelayer

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	commercelayer "github.com/incentro-dc/go-commercelayer-sdk/api"
	"net/http"
	"strings"
)

func testAccCheckDeliveryLeadTimeDestroy(s *terraform.State) error {
	client := testAccProviderCommercelayer.Meta().(*commercelayer.APIClient)

	for _, rs := range s.RootModule().Resources {
		if rs.Type == "commercelayer_delivery_lead_time" {
			_, resp, err := client.DeliveryLeadTimesApi.GETDeliveryLeadTimesDeliveryLeadTimeId(context.Background(), rs.Primary.ID).Execute()
			if resp.StatusCode == 404 {
				fmt.Printf("commercelayer_delivery_lead_time with id %s has been removed\n", rs.Primary.ID)
				continue
			}
			if err != nil {
				return err
			}

			return fmt.Errorf("received response code with status %d", resp.StatusCode)
		}

		if rs.Type == "commercelayer_address" {
			_, resp, err := client.AddressesApi.GETAddressesAddressId(context.Background(), rs.Primary.ID).Execute()
			if resp.StatusCode == 404 {
				fmt.Printf("commercelayer_address with id %s has been removed\n", rs.Primary.ID)
				continue
			}
			if err != nil {
				return err
			}

			return fmt.Errorf("received response code with status %d", resp.StatusCode)
		}

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

		for _, rs := range s.RootModule().Resources {
			if rs.Type == "commercelayer_shipping_method" {
				err := retryRemoval(10, func() (*http.Response, error) {
					_, resp, err := client.ShippingMethodsApi.GETShippingMethodsShippingMethodId(context.Background(), rs.Primary.ID).
						Execute()
					return resp, err
				})
				if err != nil {
					return err
				}
			}
		}
	}

	return nil
}

func (s *AcceptanceSuite) TestAccDeliveryLeadTime_basic() {
	resourceName := "commercelayer_delivery_lead_time.incentro_delivery_lead_time"

	resource.Test(s.T(), resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(s)
		},
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckDeliveryLeadTimeDestroy,
		Steps: []resource.TestStep{
			{
				Config: strings.Join([]string{testAccShippingMethodCreate(resourceName), testAccAddressCreate(resourceName), testAccStockLocationCreate(resourceName), testAccDeliveryLeadTimeCreate(resourceName)}, "\n"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "type", deliveryLeadTimesType),
					resource.TestCheckResourceAttr(resourceName, "attributes.0.min_hours", "10"),
					resource.TestCheckResourceAttr(resourceName, "attributes.0.metadata.foo", "bar"),
				),
			},
			{
				Config: strings.Join([]string{testAccShippingMethodCreate(resourceName), testAccAddressCreate(resourceName), testAccStockLocationCreate(resourceName), testAccDeliveryLeadTimeUpdate(resourceName)}, "\n"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "attributes.0.min_hours", "20"),
					resource.TestCheckResourceAttr(resourceName, "attributes.0.metadata.bar", "foo"),
				),
			},
		},
	})
}

func testAccDeliveryLeadTimeCreate(testName string) string {
	return hclTemplate(`
		resource "commercelayer_delivery_lead_time" "incentro_delivery_lead_time" {
		  attributes {
			min_hours = 10
			max_hours = 100
			metadata = {
			  foo : "bar"
		 	  testName: "{{.testName}}"
			}
		  }

		  relationships {
			stock_location_id = commercelayer_stock_location.incentro_stock_location.id
			shipping_method_id = commercelayer_shipping_method.incentro_shipping_method.id
		  }
		}
	`, map[string]any{"testName": testName})
}

func testAccDeliveryLeadTimeUpdate(testName string) string {
	return hclTemplate(`
		resource "commercelayer_delivery_lead_time" "incentro_delivery_lead_time" {
		  attributes {
			min_hours = 20
			max_hours = 200
			metadata = {
			  bar : "foo"
		 	  testName: "{{.testName}}"
			}
		  }

		  relationships {
			stock_location_id = commercelayer_stock_location.incentro_stock_location.id
			shipping_method_id = commercelayer_shipping_method.incentro_shipping_method.id
		  }
}
	`, map[string]any{"testName": testName})
}
