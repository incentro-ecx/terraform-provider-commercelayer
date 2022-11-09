package commercelayer

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	commercelayer "github.com/incentro-dc/go-commercelayer-sdk/api"
	"net/http"
)

func testAccCheckShippingCategoryDestroy(s *terraform.State) error {
	client := testAccProviderCommercelayer.Meta().(*commercelayer.APIClient)

	for _, rs := range s.RootModule().Resources {
		if rs.Type == "commercelayer_shipping_category" {
			err := retryRemoval(10, func() (*http.Response, error) {
				_, resp, err := client.ShippingCategoriesApi.GETShippingCategoriesShippingCategoryId(context.Background(), rs.Primary.ID).
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

func (s *AcceptanceSuite) TestAccShippingCategory_basic() {
	resourceName := "commercelayer_shipping_category.incentro_shipping_category"

	resource.Test(s.T(), resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(s)
		},
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckShippingCategoryDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccShippingCategoryCreate(resourceName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "type", shippingCategoryType),
					resource.TestCheckResourceAttr(resourceName, "attributes.0.name", "Incentro Shipping Category"),
					resource.TestCheckResourceAttr(resourceName, "attributes.0.metadata.foo", "bar"),
				),
			},
			{
				Config: testAccShippingCategoryUpdate(resourceName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "attributes.0.name", "Incentro Shipping Category Updated"),
					resource.TestCheckResourceAttr(resourceName, "attributes.0.metadata.bar", "foo"),
				),
			},
		},
	})
}

func testAccShippingCategoryCreate(testName string) string {
	return hclTemplate(`
		resource "commercelayer_shipping_category" "incentro_shipping_category" {
		  attributes {
			name                   = "Incentro Shipping Category"
			metadata               = {
			  foo : "bar"
		 	  testName: "{{.testName}}"
			}
		  }
		}
	`, map[string]any{"testName": testName})
}

func testAccShippingCategoryUpdate(testName string) string {
	return hclTemplate(`
		resource "commercelayer_shipping_category" "incentro_shipping_category" {
		  attributes {
			name                   = "Incentro Shipping Category Updated"
			metadata               = {
			  bar : "foo"
		 	  testName: "{{.testName}}"
			}
		  }
		}
	`, map[string]any{"testName": testName})
}
