package commercelayer

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	commercelayer "github.com/incentro-dc/go-commercelayer-sdk/api"
	"strings"
)

func testAccCheckBingGeocoderDestroy(s *terraform.State) error {
	client := testAccProviderCommercelayer.Meta().(*commercelayer.APIClient)

	for _, rs := range s.RootModule().Resources {
		if rs.Type == "commercelayer_bing_geocoder" {
			_, resp, err := client.BingGeocodersApi.GETBingGeocodersBingGeocoderId(context.Background(), rs.Primary.ID).Execute()
			if resp.StatusCode == 404 {
				fmt.Printf("commercelayer_bing_geocoder with id %s has been removed\n", rs.Primary.ID)
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

func (s *AcceptanceSuite) TestAccBingGeocoder_basic() {
	resourceName := "commercelayer_bing_geocoder.incentro_bing_geocoder"

	resource.Test(s.T(), resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(s)
		},
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckBingGeocoderDestroy,
		Steps: []resource.TestStep{
			{
				Config: strings.Join([]string{testAccBingGeocoderCreate(resourceName)}, "\n"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "type", bingGeocodersType),
					resource.TestCheckResourceAttr(resourceName, "attributes.0.name", "Incentro Bing Geocoder"),
					resource.TestCheckResourceAttr(resourceName, "attributes.0.metadata.foo", "bar"),
				),
			},
			{
				Config: strings.Join([]string{testAccBingGeocoderUpdate(resourceName)}, "\n"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "attributes.0.name", "Incentro Updated Bing Geocoder"),
					resource.TestCheckResourceAttr(resourceName, "attributes.0.metadata.bar", "foo"),
				),
			},
		},
	})
}

func testAccBingGeocoderCreate(testName string) string {
	return hclTemplate(`
		resource "commercelayer_bing_geocoder" "incentro_bing_geocoder" {
  			attributes {
    			name                   = "Incentro Bing Geocoder"
    			key               	   = "Bing Virtualearth Key"
				metadata = {
			  		foo : "bar"
		 	 		testName: "{{.testName}}"
				}
  			}
	}`, map[string]any{"testName": testName})
}

func testAccBingGeocoderUpdate(testName string) string {
	return hclTemplate(`
		resource "commercelayer_bing_geocoder" "incentro_bing_geocoder" {
  			attributes {
    			name                   = "Incentro Updated Bing Geocoder"
    			key                    = "Bing Virtualearth Key"
				metadata = {
			  		bar : "foo"
		 	 		testName: "{{.testName}}"
				}
  			}
	}`, map[string]any{"testName": testName})
}
