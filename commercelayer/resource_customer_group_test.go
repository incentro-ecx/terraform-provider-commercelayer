package commercelayer

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	commercelayer "github.com/incentro-dc/go-commercelayer-sdk/api"
)

func testAccCheckCustomerGroupDestroy(s *terraform.State) error {
	client := testAccProviderCommercelayer.Meta().(*commercelayer.APIClient)

	for _, rs := range s.RootModule().Resources {
		if rs.Type == "commercelayer_customer_group" {
			_, resp, err := client.CustomerGroupsApi.GETCustomerGroupsCustomerGroupId(context.Background(), rs.Primary.ID).Execute()
			if resp.StatusCode == 404 {
				fmt.Printf("commercelayer_customer_group with id %s has been removed\n", rs.Primary.ID)
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

func (s *AcceptanceSuite) TestAccCustomerGroup_basic() {
	resourceName := "commercelayer_customer_group.incentro_customer_group"

	resource.Test(s.T(), resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(s)
		},
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckCustomerGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCustomerGroupCreate(resourceName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "type", customerGroupType),
					resource.TestCheckResourceAttr(resourceName, "attributes.0.name", "Incentro customer group"),
					resource.TestCheckResourceAttr(resourceName, "attributes.0.metadata.foo", "bar"),
				),
			},
			{
				Config: testAccCustomerGroupUpdate(resourceName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "attributes.0.name", "Incentro updated customer group"),
					resource.TestCheckResourceAttr(resourceName, "attributes.0.metadata.bar", "foo"),
				),
			},
		},
	})
}

func testAccCustomerGroupCreate(testName string) string {
	return hclTemplate(`
		resource "commercelayer_customer_group" "incentro_customer_group" {
		  attributes {
			name = "Incentro customer group"
			metadata = {
			  foo : "bar"
			  testName: "{{.testName}}"
			}
		  }
		}
	`, map[string]any{"testName": testName})
}

func testAccCustomerGroupUpdate(testName string) string {
	return hclTemplate(`
		resource "commercelayer_customer_group" "incentro_customer_group" {
		  attributes {
			name = "Incentro updated customer group"
			metadata = {
			  bar : "foo"
			  testName: "{{.testName}}"
			}
		  }
		}
	`, map[string]any{"testName": testName})
}
