package commercelayer

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	commercelayer "github.com/incentro-dc/go-commercelayer-sdk/api"
	"regexp"
	"testing"
)

func testAccCheckPriceListDestroy(s *terraform.State) error {
	client := testAccProviderCommercelayer.Meta().(*commercelayer.APIClient)

	for _, rs := range s.RootModule().Resources {
		if rs.Type == "commercelayer_price_list" {
			_, resp, err := client.PriceListsApi.GETPriceListsPriceListId(context.Background(), rs.Primary.ID).Execute()
			if resp.StatusCode == 404 {
				fmt.Printf("commercelayer_price_list with id %s has been removed\n", rs.Primary.ID)
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

func TestAccPriceList_basic(t *testing.T) {
	resourceName := "commercelayer_price_list.incentro_price_list"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckPriceListDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccPriceListCreate(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "attributes.0.name", "incentro price list"),
					resource.TestCheckResourceAttr(resourceName, "attributes.0.currency_code", "EUR"),
					resource.TestCheckResourceAttr(resourceName, "attributes.0.metadata.foo", "bar"),
				),
			},
			{
				Config: testAccPriceListUpdate(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "attributes.0.name", "incentro updated price list"),
					resource.TestCheckResourceAttr(resourceName, "attributes.0.currency_code", "CHF"),
					resource.TestCheckResourceAttr(resourceName, "attributes.0.metadata.bar", "foo"),
				),
			},
		},
	})
}

func TestAccPriceList_invalid(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckPriceListDestroy,
		Steps: []resource.TestStep{
			{
				Config:      testAccPriceListCreateInvalidCurrency(),
				ExpectError: regexp.MustCompile(".*FOOBAR.*"),
			},
		},
	})
}

func testAccPriceListCreate() string {
	return hclTemplate(`
		resource "commercelayer_price_list" "incentro_price_list" {
		  attributes {
			name          = "incentro price list"
			currency_code = "EUR"
			metadata = {
			  foo : "bar"
			}
		  }
		}
	`, map[string]any{})
}

func testAccPriceListUpdate() string {
	return hclTemplate(`
		resource "commercelayer_price_list" "incentro_price_list" {
		  attributes {
			name          = "incentro updated price list"
			currency_code = "CHF"
			metadata = {
			  bar : "foo"
			}
		  }
		}
	`, map[string]any{})
}

func testAccPriceListCreateInvalidCurrency() string {
	return hclTemplate(`
		resource "commercelayer_price_list" "incentro_price_list_invalid_currency" {
		  attributes {
			name          = "incentro updated price list"
			currency_code = "FOOBAR"
			metadata = {
			  bar : "foo"
			}
		  }
		}
	`, map[string]any{})
}
