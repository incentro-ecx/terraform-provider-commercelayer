package commercelayer

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	commercelayer "github.com/incentro-dc/go-commercelayer-sdk/api"
)

func testAccCheckTaxjarAccountDestroy(s *terraform.State) error {
	client := testAccProviderCommercelayer.Meta().(*commercelayer.APIClient)

	for _, rs := range s.RootModule().Resources {
		if rs.Type == "commercelayer_taxjar_accounts" {
			_, resp, err := client.TaxjarAccountsApi.
				GETTaxjarAccountsTaxjarAccountId(context.Background(), rs.Primary.ID).Execute()
			if resp.StatusCode == 404 {
				fmt.Printf("commercelayer_taxjar_accounts with id %s has been removed\n", rs.Primary.ID)
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

func (s *AcceptanceSuite) TestAccTaxjarAccount_basic() {
	resourceName := "commercelayer_taxjar_accounts.incentro_taxjar_account"

	resource.Test(s.T(), resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(s)
		},
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckTaxjarAccountDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTaxjarAccountCreate(resourceName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "type", taxjarAccountsType),
					resource.TestCheckResourceAttr(resourceName, "attributes.0.name", "Incentro Taxjar Account"),
					resource.TestCheckResourceAttr(resourceName, "attributes.0.metadata.foo", "bar"),
				),
			},
			{
				Config: testAccTaxjarAccountUpdate(resourceName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "attributes.0.name", "Incentro Taxjar Account Changed"),
					resource.TestCheckResourceAttr(resourceName, "attributes.0.metadata.bar", "foo"),
				),
			},
		},
	})
}

func testAccTaxjarAccountCreate(testName string) string {
	return hclTemplate(`
		resource "commercelayer_taxjar_accounts" "incentro_taxjar_account" {
           attributes {
			name = "Incentro Taxjar Account"
			api_key = "TAXJAR_API_KEY"
			metadata = {
				foo: "bar"
				testName: "{{.testName}}"
    		}
  		}
	}
`, map[string]any{"testName": testName})
}

func testAccTaxjarAccountUpdate(testName string) string {
	return hclTemplate(`
		resource "commercelayer_taxjar_accounts" "incentro_taxjar_account" {
           attributes {
			name                   = "Incentro Taxjar Account Changed"
			api_key = "TAXJAR_API_KEY"
			metadata = {
				bar: "foo"
				testName: "{{.testName}}"
    		}
  		}
	}
`, map[string]any{"testName": testName})
}
