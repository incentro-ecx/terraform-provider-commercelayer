package commercelayer

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	commercelayer "github.com/incentro-dc/go-commercelayer-sdk/api"
	"strings"
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

func (s *AcceptanceSuite) TestAccMerchant_basic() {
	resourceName := "commercelayer_merchant.incentro_merchant"

	resource.Test(s.T(), resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(s)
		},
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckMerchantDestroy,
		Steps: []resource.TestStep{
			{
				Config: strings.Join([]string{testAccAddressCreate(resourceName), testAccMerchantCreate(resourceName)}, "\n"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "type", merchantType),
					resource.TestCheckResourceAttr(resourceName, "attributes.0.name", "Incentro Merchant"),
					resource.TestCheckResourceAttr(resourceName, "attributes.0.metadata.foo", "bar"),
				),
			},
			{
				Config: strings.Join([]string{testAccAddressCreate(resourceName), testAccMerchantUpdate(resourceName)}, "\n"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "attributes.0.name", "Incentro Updated Merchant"),
					resource.TestCheckResourceAttr(resourceName, "attributes.0.metadata.bar", "foo"),
				),
			},
		},
	})
}

func testAccMerchantCreate(testName string) string {
	return hclTemplate(`
		resource "commercelayer_merchant" "incentro_merchant" {
		  attributes {
			name = "Incentro Merchant"
			metadata = {
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

func testAccMerchantUpdate(testName string) string {
	return hclTemplate(`
		resource "commercelayer_merchant" "incentro_merchant" {
		  attributes {
			name = "Incentro Updated Merchant"
			metadata = {
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
