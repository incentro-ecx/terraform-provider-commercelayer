package commercelayer

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	commercelayer "github.com/incentro-dc/go-commercelayer-sdk/api"
	"net/http"
	"strings"
)

func testAccCheckSkuDestroy(s *terraform.State) error {
	client := testAccProviderCommercelayer.Meta().(*commercelayer.APIClient)

	for _, rs := range s.RootModule().Resources {
		if rs.Type == "commercelayer_sku" {
			err := retryRemoval(10, func() (*http.Response, error) {
				_, resp, err := client.SkusApi.GETSkusSkuId(context.Background(), rs.Primary.ID).
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

func (s *AcceptanceSuite) TestAccSku_basic() {
	resourceName := "commercelayer_sku.incentro_sku"

	resource.Test(s.T(), resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(s)
		},
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckSkuDestroy,
		Steps: []resource.TestStep{
			{
				Config: strings.Join([]string{testAccSkuCreate(resourceName), testAccShippingCategoryCreate(resourceName)}, "\n"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "type", skusType),
					resource.TestCheckResourceAttr(resourceName, "attributes.0.name", "Incentro SKU"),
					resource.TestCheckResourceAttr(resourceName, "attributes.0.code", "TSHIRTMM000000FFFFFFXLXX"),
					resource.TestCheckResourceAttr(resourceName, "attributes.0.metadata.foo", "bar"),
				),
			},
			{
				Config: strings.Join([]string{testAccSkuUpdate(resourceName), testAccShippingCategoryCreate(resourceName)}, "\n"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "attributes.0.name", "Incentro SKU Updated"),
					resource.TestCheckResourceAttr(resourceName, "attributes.0.code", "TSHIRTMM000000FFFFFFXLXX"),
					resource.TestCheckResourceAttr(resourceName, "attributes.0.metadata.bar", "foo"),
				),
			},
		},
	})
}

func testAccSkuCreate(testName string) string {
	return hclTemplate(`
		resource "commercelayer_sku" "incentro_sku" {
		  attributes {
			name = "Incentro SKU"
			code = "TSHIRTMM000000FFFFFFXLXX"
			metadata = {
			  foo : "bar"
			  testName : "{{.testName}}"
			}
		  }
		  relationships {
			shipping_category_id = commercelayer_shipping_category.incentro_shipping_category.id
		  }
		}
	`, map[string]any{"testName": testName})
}

func testAccSkuUpdate(testName string) string {
	return hclTemplate(`
		resource "commercelayer_sku" "incentro_sku" {
		  attributes {
			name = "Incentro SKU Updated"
			code = "TSHIRTMM000000FFFFFFXLXX"
			metadata = {
			  bar : "foo"
			  testName : "{{.testName}}"
			}
		  }
		  relationships {
			shipping_category_id = commercelayer_shipping_category.incentro_shipping_category.id
		  }
		}
	`, map[string]any{"testName": testName})
}
