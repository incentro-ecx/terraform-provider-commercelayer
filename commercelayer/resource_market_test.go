package commercelayer

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	commercelayer "github.com/incentro-dc/go-commercelayer-sdk/api"
	"net/http"
	"strings"
	"testing"
)

func testAccCheckMarketDestroy(s *terraform.State) error {
	client := testAccProviderCommercelayer.Meta().(*commercelayer.APIClient)

	for _, rs := range s.RootModule().Resources {
		if rs.Type == "commercelayer_market" {
			err := retryRemoval(10, func() (*http.Response, error) {
				_, resp, err := client.MarketsApi.
					GETMarketsMarketId(context.Background(), rs.Primary.ID).
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

func TestAccMarket_basic(t *testing.T) {
	resourceName := "commercelayer_market.incentro_market"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckMarketDestroy,
		Steps: []resource.TestStep{
			{
				Config: strings.Join([]string{
					testAccAddressCreate(),
					testAccInventoryModelCreate(),
					testAccMerchantCreate(),
					testAccPriceListCreate(),
					testAccMarketCreate()}, "\n",
				),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "attributes.0.name", "Incentro Market"),
					resource.TestCheckResourceAttr(resourceName, "attributes.0.facebook_pixel_id", "pixel"),
				),
			},
			{
				Config: strings.Join([]string{
					testAccAddressCreate(),
					testAccInventoryModelCreate(),
					testAccMerchantCreate(),
					testAccPriceListCreate(),
					testAccMarketUpdate()}, "\n",
				),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "attributes.0.name", "Incentro Market Changed"),
					resource.TestCheckResourceAttr(resourceName, "attributes.0.facebook_pixel_id", "pixelchanged"),
				),
			},
		},
	})
}

func testAccMarketCreate() string {
	return hclTemplate(`
		resource "commercelayer_market" "incentro_market" {
		  attributes {
			name              = "Incentro Market"
			facebook_pixel_id = "pixel"
		  }
		
		  relationships {
			inventory_model_id = commercelayer_inventory_model.incentro_inventory_model.id
			merchant_id        = commercelayer_merchant.incentro_merchant.id
			price_list_id      = commercelayer_price_list.incentro_price_list.id
		  }
		}
	`, map[string]any{})
}

func testAccMarketUpdate() string {
	return hclTemplate(`
		resource "commercelayer_market" "incentro_market" {
		  attributes {
			name              = "Incentro Market Changed"
			facebook_pixel_id = "pixelchanged"
		  }
		
		  relationships {
			inventory_model_id = commercelayer_inventory_model.incentro_inventory_model.id
			merchant_id        = commercelayer_merchant.incentro_merchant.id
			price_list_id      = commercelayer_price_list.incentro_price_list.id
		  }
		}
	`, map[string]any{})
}
