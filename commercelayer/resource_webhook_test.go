package commercelayer

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	commercelayer "github.com/incentro-dc/go-commercelayer-sdk/api"
	"testing"
)

func testAccCheckWebhookDestroy(s *terraform.State) error {
	client := testAccProviderCommercelayer.Meta().(*commercelayer.APIClient)

	for _, rs := range s.RootModule().Resources {
		if rs.Type == "commercelayer_webhook" {
			_, resp, err := client.WebhooksApi.GETWebhooksWebhookId(context.Background(), rs.Primary.ID).Execute()
			if resp.StatusCode == 404 {
				fmt.Printf("commercelayer_webhook with id %s has been removed\n", rs.Primary.ID)
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

func TestAccWebhook_basic(t *testing.T) {
	resourceName := "commercelayer_webhook.incentro_webhook"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckWebhookDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccWebhookCreate(resourceName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "type", webhookType),
					resource.TestCheckResourceAttr(resourceName, "attributes.0.name", "incentro webhook"),
					resource.TestCheckResourceAttr(resourceName, "attributes.0.topic", "orders.create"),
					resource.TestCheckResourceAttr(resourceName, "attributes.0.callback_url", "http://example.url"),
					resource.TestCheckResourceAttr(resourceName, "attributes.0.include_resources.0", "customer"),
					resource.TestCheckResourceAttr(resourceName, "attributes.0.metadata.foo", "bar"),
				),
			},
			{
				Config: testAccWebhookUpdate(resourceName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "attributes.0.name", "incentro updated webhook"),
					resource.TestCheckResourceAttr(resourceName, "attributes.0.topic", "orders.place"),
					resource.TestCheckResourceAttr(resourceName, "attributes.0.callback_url", "http://other-example.url"),
					resource.TestCheckResourceAttr(resourceName, "attributes.0.include_resources.0", "line_items"),
					resource.TestCheckResourceAttr(resourceName, "attributes.0.metadata.bar", "foo"),
				),
			},
		},
	})
}

func testAccWebhookCreate(testName string) string {
	return hclTemplate(`
		resource "commercelayer_webhook" "incentro_webhook" {
		  attributes {
			name         = "incentro webhook"
			topic        = "orders.create"
			callback_url = "http://example.url"
			include_resources = [
			  "customer"
			]
			metadata = {
			  foo : "bar"
		 	  testName: "{{.testName}}"
			}
		  }
		}
	`, map[string]any{"testName": testName})
}

func testAccWebhookUpdate(testName string) string {
	return hclTemplate(`
		resource "commercelayer_webhook" "incentro_webhook" {
		  attributes {
			name         = "incentro updated webhook"
			topic        = "orders.place"
			callback_url = "http://other-example.url"
			include_resources = [
			  "line_items"
			]
			metadata = {
			  bar : "foo"
		 	  testName: "{{.testName}}"
			}
		  }
		}
	`, map[string]any{"testName": testName})
}
