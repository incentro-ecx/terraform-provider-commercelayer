package commercelayer

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	commercelayer "github.com/incentro-dc/go-commercelayer-sdk/api"
	"strings"
)

func testAccCheckDeliveryLeadTimeDestroy(s *terraform.State) error {
	client := testAccProviderCommercelayer.Meta().(*commercelayer.APIClient)

	for _, rs := range s.RootModule().Resources {
		if rs.Type == "commercelayer_delivery_lead_type" {
			_, resp, err := client.DeliveryLeadTimesApi.GETDeliveryLeadTimesDeliveryLeadTimeId(context.Background(), rs.Primary.ID).Execute()
			if resp.StatusCode == 404 {
				fmt.Printf("commercelayer_merchant with id %s has been removed\n", rs.Primary.ID)
				continue
			}
			if err != nil {
				return err
			}

			return fmt.Errorf("received response code with status %d", resp.StatusCode)
		}

		if rs.Type == "commercelayer_stock_location" {
			_, resp, err := client.StockLocationsApi.GETStockLocationsStockLocationId(context.Background(), rs.Primary.ID).Execute()
			if resp.StatusCode == 404 {
				fmt.Printf("commercelayer_address with id %s has been removed\n", rs.Primary.ID)
				continue
			}
			if err != nil {
				return err
			}

			return fmt.Errorf("received response code with status %d", resp.StatusCode)
		}

		if rs.Type == "commercelayer_shipping_method" {
			_, resp, err := client.ShippingMethodsApi.GETShippingMethodsShippingMethodId(context.Background(), rs.Primary.ID).Execute()
			if resp.StatusCode == 404 {
				fmt.Printf("commercelayer_address with id %s has been removed\n", rs.Primary.ID)
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
				Config: strings.Join([]string{testAccShippingMethodCreate(resourceName), testAccStockLocationCreate(resourceName), testAccDeliveryLeadTimeCreate(resourceName)}, "\n"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "type", deliveryLeadTimesType),
					resource.TestCheckResourceAttr(resourceName, "attributes.0.min_hours", "10"),
					resource.TestCheckResourceAttr(resourceName, "attributes.0.metadata.foo", "bar"),
				),
			},
			{
				Config: strings.Join([]string{testAccShippingMethodCreate(resourceName), testAccStockLocationCreate(resourceName), testAccDeliveryLeadTimeUpdate(resourceName)}, "\n"),
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
			}
		  }
		
		  relationships {
			stock_location = commercelayer_stock_location.incentro_warehouse_location.id
			shipping_method = commercelayer_shipping_method.incentro_shipping_method.id
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
			  foo : "bar"
			}
		  }
		
		  relationships {
			stock_location = commercelayer_stock_location.incentro_warehouse_location.id
			shipping_method = commercelayer_shipping_method.incentro_shipping_method.id
		  }
}
	`, map[string]any{"testName": testName})
}
