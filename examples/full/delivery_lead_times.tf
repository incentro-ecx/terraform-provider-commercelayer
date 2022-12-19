resource "commercelayer_delivery_lead_time" "incentro_delivery_lead_time" {
  attributes {
    min_hours = 10
    max_hours = 100
    metadata = {
      foo : "bar"
    }
  }

  relationships {
    stock_location = commercelayer_stock_location.incentro_warehouse_location.id
    shipping_method = commercelayer_shipping_method.incentro_shipping_method.id
  }
}