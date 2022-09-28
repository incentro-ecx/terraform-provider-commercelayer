package commercelayer

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	commercelayer "github.com/incentro-dc/go-commercelayer-sdk/api"
	"testing"
)

func testAccCheckMerchantDestroy(s *terraform.State) error {
	client := testAccProviderCommercelayer.Meta().(*commercelayer.APIClient)

	for _, rs := range s.RootModule().Resources {
		if rs.Type == "commercelayer_merchant" {
			_, resp, err := client.MerchantsApi.GETMerchantsMerchantId(context.Background(), rs.Primary.ID).Execute()
			if resp.StatusCode == 404 {
				fmt.Printf("commercelayer_merchant with id %s has been removed\n", rs.Primary.ID)
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

	}
	return nil
}

func TestAccMerchant_basic(t *testing.T) {
	resourceName := "commercelayer_merchant.incentro_merchant"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckMerchantDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMerchantCreate(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "attributes.0.name", "Incentro Merchant"),
					resource.TestCheckResourceAttr(resourceName, "attributes.0.metadata.foo", "bar"),
				),
			},
			{
				Config: testAccMerchantUpdate(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "attributes.0.name", "Incentro Updated Merchant"),
					resource.TestCheckResourceAttr(resourceName, "attributes.0.metadata.bar", "foo"),
				),
			},
		},
	})
}

func testAccMerchantCreate() string {
	return hclTemplate(`
		resource "commercelayer_address" "incentro_address" {
		  attributes {
			business     = true
			company      = "Incentro"
			line_1       = "Van Nelleweg 1"
			zip_code     = "3044 BC"
			country_code = "NL"
			city         = "Rotterdam"
			phone        = "+31(0)10 20 20 544"
			state_code   = "ZH"
			metadata = {
			  foo : "bar"
			}
		  }
		}

		resource "commercelayer_merchant" "incentro_merchant" {
		  attributes {
			name = "Incentro Merchant"
			metadata = {
			  foo : "bar"
			}
		  }
		
		  relationships {
			address = commercelayer_address.incentro_address.id
		  }
		}
	`, map[string]any{})
}

func testAccMerchantUpdate() string {
	return hclTemplate(`
		resource "commercelayer_address" "incentro_address" {
		  attributes {
			business     = true
			company      = "Incentro"
			line_1       = "Van Nelleweg 1"
			zip_code     = "3044 BC"
			country_code = "NL"
			city         = "Rotterdam"
			phone        = "+31(0)10 20 20 544"
			state_code   = "ZH"
			metadata = {
			  foo : "bar"
			}
		  }
		}

		resource "commercelayer_merchant" "incentro_merchant" {
		  attributes {
			name = "Incentro Updated Merchant"
			metadata = {
			  bar : "foo"
			}
		  }
		
		  relationships {
			address = commercelayer_address.incentro_address.id
		  }
		}
	`, map[string]any{})
}
