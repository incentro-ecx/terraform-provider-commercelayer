package commercelayer

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	commercelayer "github.com/incentro-dc/go-commercelayer-sdk/api"
	"testing"
)

func testAccCheckExternalTaxCalculatorDestroy(s *terraform.State) error {
	client := testAccProviderCommercelayer.Meta().(*commercelayer.APIClient)

	for _, rs := range s.RootModule().Resources {
		if rs.Type == "commercelayer_external_tax_calculator" {
			_, resp, err := client.ExternalTaxCalculatorsApi.GETExternalTaxCalculatorsExternalTaxCalculatorId(context.Background(), rs.Primary.ID).Execute()
			if resp.StatusCode == 404 {
				fmt.Printf("commercelayer_external_tax_calculator with id %s has been removed\n", rs.Primary.ID)
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

func TestAccExternalTaxCalculator_basic(t *testing.T) {
	resourceName := "commercelayer_external_tax_calculator.incentro_external_tax_calculator"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckExternalTaxCalculatorDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccExternalTaxCalculatorCreate(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "type", externalTaxCalculatorType),
					resource.TestCheckResourceAttr(resourceName, "attributes.0.name", "incentro_external_tax_calculator"),
					resource.TestCheckResourceAttr(resourceName, "attributes.0.metadata.foo", "bar"),
					resource.TestCheckResourceAttr(resourceName, "attributes.0.tax_calculator_url", "https://example.com"),
				),
			},
			{
				Config: testAccExternalTaxCalculatorUpdateWithUrls(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "attributes.0.name", "incentro_external_tax_calculator_changed"),
					resource.TestCheckResourceAttr(resourceName, "attributes.0.metadata.bar", "foo"),
					resource.TestCheckResourceAttr(resourceName, "attributes.0.tax_calculator_url", "https://foo.com"),
				),
			},
		},
	})
}

func testAccExternalTaxCalculatorCreate() string {
	return hclTemplate(`
	resource "commercelayer_external_tax_calculator" "incentro_external_tax_calculator" {
	  attributes {
		name          = "incentro_external_tax_calculator"
		tax_calculator_url = "https://example.com"
		metadata = {
		  foo : "bar"
		}
	  }
	}
	`, map[string]any{})
}

func testAccExternalTaxCalculatorUpdateWithUrls() string {
	return hclTemplate(`
		resource "commercelayer_external_tax_calculator" "incentro_external_tax_calculator" {
		  attributes {
			name          = "incentro_external_tax_calculator_changed"
			tax_calculator_url = "https://foo.com"
			metadata = {
			  bar : "foo"
			}
		  }
		}
	`, map[string]any{})
}
