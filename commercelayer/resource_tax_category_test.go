package commercelayer

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	commercelayer "github.com/incentro-dc/go-commercelayer-sdk/api"
	"strings"
)

func testAccCheckTaxCategoryDestroy(s *terraform.State) error {
	client := testAccProviderCommercelayer.Meta().(*commercelayer.APIClient)

	for _, rs := range s.RootModule().Resources {
		if rs.Type == "commercelayer_tax_category" {
			_, resp, err := client.TaxCategoriesApi.GETTaxCategoriesTaxCategoryId(context.Background(), rs.Primary.ID).Execute()
			if resp.StatusCode == 404 {
				fmt.Printf("commercelayer_tax_category with id %s has been removed\n", rs.Primary.ID)
				continue
			}
			if err != nil {
				return err
			}

			return fmt.Errorf("received response code with status %d", resp.StatusCode)
		}

		if rs.Type == "commercelayer_manual_tax_calculator" {
			_, resp, err := client.ManualTaxCalculatorsApi.GETManualTaxCalculatorsManualTaxCalculatorId(context.Background(), rs.Primary.ID).Execute()
			if resp.StatusCode == 404 {
				fmt.Printf("commercelayer_manual_tax_calculator with id %s has been removed\n", rs.Primary.ID)
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

func (s *AcceptanceSuite) TestAccTaxCategory_basic() {
	resourceName := "commercelayer_tax_category.incentro_tax_category"

	resource.Test(s.T(), resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(s)
		},
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckTaxCategoryDestroy,
		Steps: []resource.TestStep{
			{
				Config: strings.Join([]string{testAccTaxCategoryCreate(resourceName), testAccManualTaxCalculatorCreate(resourceName)}, "\n"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "type", taxCategoriesType),
					resource.TestCheckResourceAttr(resourceName, "attributes.0.metadata.foo", "bar"),
				),
			},
			{
				Config: strings.Join([]string{testAccTaxCategoryUpdate(resourceName), testAccManualTaxCalculatorCreate(resourceName)}, "\n"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "attributes.0.metadata.bar", "foo"),
				),
			},
		},
	})
}

// TODO: check if sku_id is correct like this
func testAccTaxCategoryCreate(testName string) string {
	return hclTemplate(`
	resource "commercelayer_tax_category" "incentro_tax_category" {
	  attributes {
		code          = "31000"
		metadata = {
		  foo : "bar"
		  testName: "{{.testName}}"
		}
	  }
	  relationships {
		sku_id = "TSHIRTWF000000E63E74MXXX"
        tax_calculator_id = commercelayer_manual_tax_calculator.incentro_manual_tax_calculator.id
      }
	}
	`, map[string]any{"testName": testName})
}

// TODO: check if sku_id is correct like this
func testAccTaxCategoryUpdate(testName string) string {
	return hclTemplate(`
		resource "commercelayer_tax_category" "incentro_tax_category" {
		  attributes {
		    code          = "31000"
			metadata = {
			  bar : "foo"
		 	  testName: "{{.testName}}"
			}
		  }
		  relationships {
			sku_id = "TSHIRTWF000000E63E74MXXX"
			}
		}
	`, map[string]any{"testName": testName})
}
