package commercelayer

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	commercelayer "github.com/incentro-dc/go-commercelayer-sdk/api"
	"testing"
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

func TestAccCustomerGroup_basic(t *testing.T) {
	resourceName := "commercelayer_customer_group.incentro_customer_group"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckCustomerGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCustomerGroupCreate(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "attributes.0.name", "Incentro customer group"),
					resource.TestCheckResourceAttr(resourceName, "attributes.0.metadata.foo", "bar"),
				),
			},
			{
				Config: testAccCustomerGroupUpdate(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "attributes.0.name", "Incentro updated customer group"),
					resource.TestCheckResourceAttr(resourceName, "attributes.0.metadata.bar", "foo"),
				),
			},
		},
	})
}

func testAccCustomerGroupCreate() string {
	return hclTemplate(`
		resource "commercelayer_customer_group" "incentro_customer_group" {
		  attributes {
			name = "Incentro customer group"
			metadata = {
			  foo : "bar"
			}
		  }
		}
	`, map[string]any{})
}

func testAccCustomerGroupUpdate() string {
	return hclTemplate(`
		resource "commercelayer_customer_group" "incentro_customer_group" {
		  attributes {
			name = "Incentro updated customer group"
			metadata = {
			  bar : "foo"
			}
		  }
		}
	`, map[string]any{})
}
