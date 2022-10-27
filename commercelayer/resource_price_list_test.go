package commercelayer

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	commercelayer "github.com/incentro-dc/go-commercelayer-sdk/api"
	"net/http"
	"regexp"
)

func testAccCheckPriceListDestroy(s *terraform.State) error {
	client := testAccProviderCommercelayer.Meta().(*commercelayer.APIClient)

	for _, rs := range s.RootModule().Resources {
		if rs.Type == "commercelayer_price_list" {
			err := retryRemoval(10, func() (*http.Response, error) {
				_, resp, err := client.PriceListsApi.GETPriceListsPriceListId(context.Background(), rs.Primary.ID).
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

func (s *AcceptanceSuite) TestAccPriceList_basic() {
	resourceName := "commercelayer_price_list.incentro_price_list"

	resource.Test(s.T(), resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(s)
		},
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckPriceListDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccPriceListCreate(resourceName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "type", priceListType),
					resource.TestCheckResourceAttr(resourceName, "attributes.0.name", "incentro price list"),
					resource.TestCheckResourceAttr(resourceName, "attributes.0.currency_code", "EUR"),
					resource.TestCheckResourceAttr(resourceName, "attributes.0.metadata.foo", "bar"),
				),
			},
			{
				Config: testAccPriceListUpdate(resourceName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "attributes.0.name", "incentro updated price list"),
					resource.TestCheckResourceAttr(resourceName, "attributes.0.currency_code", "CHF"),
					resource.TestCheckResourceAttr(resourceName, "attributes.0.metadata.bar", "foo"),
				),
			},
		},
	})
}

func (s *AcceptanceSuite) TestAccPriceList_invalid() {
	resourceName := "commercelayer_price_list.incentro_price_list_invalid_currency"

	resource.Test(s.T(), resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(s)
		},
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckPriceListDestroy,
		Steps: []resource.TestStep{
			{
				Config:      testAccPriceListCreateInvalidCurrency(resourceName),
				ExpectError: regexp.MustCompile(".*FOOBAR.*"),
			},
		},
	})
}

func testAccPriceListCreate(testName string) string {
	return hclTemplate(`
		resource "commercelayer_price_list" "incentro_price_list" {
		  attributes {
			name          = "incentro price list"
			currency_code = "EUR"
			metadata = {
			  foo : "bar"
		 	  testName: "{{.testName}}"
			}
		  }
		}
	`, map[string]any{"testName": testName})
}

func testAccPriceListUpdate(testName string) string {
	return hclTemplate(`
		resource "commercelayer_price_list" "incentro_price_list" {
		  attributes {
			name          = "incentro updated price list"
			currency_code = "CHF"
			metadata = {
			  bar : "foo"
		 	  testName: "{{.testName}}"
			}
		  }
		}
	`, map[string]any{"testName": testName})
}

func testAccPriceListCreateInvalidCurrency(testName string) string {
	return hclTemplate(`
		resource "commercelayer_price_list" "incentro_price_list_invalid_currency" {
		  attributes {
			name          = "incentro updated price list"
			currency_code = "FOOBAR"
			metadata = {
			  bar : "foo"
		 	  testName: "{{.testName}}"
			}
		  }
		}
	`, map[string]any{"testName": testName})
}
