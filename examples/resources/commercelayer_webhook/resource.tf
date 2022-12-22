resource "commercelayer_webhook" "incentro_webhook" {
  attributes {
    name         = "Incentro Webhook"
    topic        = "orders.create"
    callback_url = "http://example.url"
    include_resources = [
      "customer",
      "line_items"
    ]
  }
}
