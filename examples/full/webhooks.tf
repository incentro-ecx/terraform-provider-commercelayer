resource "commercelayer_webhook" "orders_create_webhook" {
  attributes {
    name         = "order-fulfillment"
    topic        = "orders.create"
    callback_url = "http://example.url"
    include_resources = [
      "customer",
      "line_items"
    ]
    metadata = {
      foo : "bar"
    }
  }
}
