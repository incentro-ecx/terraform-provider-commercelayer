package commercelayer

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	commercelayer "github.com/incentro-dc/go-commercelayer-sdk/api"
	"strings"
)

func testAccCheckTaxRuleDestroy(s *terraform.State) error {
	client := testAccProviderCommercelayer.Meta().(*commercelayer.APIClient)

	for _, rs := range s.RootModule().Resources {
		if rs.Type == "commercelayer_tax_rule" {
			_, resp, err := client.TaxRulesApi.GETTaxRulesTaxRuleId(context.Background(), rs.Primary.ID).Execute()
			if resp.StatusCode == 404 {
				fmt.Printf("commercelayer_tax_rule with id %s has been removed\n", rs.Primary.ID)
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

func (s *AcceptanceSuite) TestAccTaxRule_basic() {
	resourceName := "commercelayer_tax_rule.incentro_tax_rule"

	resource.Test(s.T(), resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(s)
		},
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckTaxRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: strings.Join([]string{testAccTaxRuleCreate(resourceName), testAccManualTaxCalculatorCreate(resourceName)}, "\n"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "type", taxRulesType),
					resource.TestCheckResourceAttr(resourceName, "attributes.0.name", "Incentro Tax Rule"),
					resource.TestCheckResourceAttr(resourceName, "attributes.0.metadata.foo", "bar"),
				),
			},
			{
				Config: strings.Join([]string{testAccTaxRuleUpdate(resourceName), testAccManualTaxCalculatorCreate(resourceName)}, "\n"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "attributes.0.name", "Incentro Tax Rule Changed"),
					resource.TestCheckResourceAttr(resourceName, "attributes.0.metadata.bar", "foo"),
				),
			},
		},
	})
}

func testAccTaxRuleCreate(testName string) string {
	return hclTemplate(`
		resource "commercelayer_tax_rule" "incentro_tax_rule" {
		  attributes {
			name = "Incentro Tax Rule"
			metadata = {
			  foo : "bar"
			  testName : "{{.testName}}"
			}
		  }
		  relationships {
			manual_tax_calculator_id = commercelayer_manual_tax_calculator.incentro_manual_tax_calculator.id
		  }
		}
	`, map[string]any{"testName": testName})
}

func testAccTaxRuleUpdate(testName string) string {
	return hclTemplate(`
		resource "commercelayer_tax_rule" "incentro_tax_rule" {
		  attributes {
			name = "Incentro Tax Rule Changed"
			metadata = {
			  bar : "foo"
			  testName : "{{.testName}}"
			}
		  }
		  relationships {
			manual_tax_calculator_id = commercelayer_manual_tax_calculator.incentro_manual_tax_calculator.id
		  }
		}
	`, map[string]any{"testName": testName})
}
