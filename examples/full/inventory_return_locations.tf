resource "commercelayer_inventory_return_location" "incentro_return_location" {
  attributes {
    priority = 1
  }

  relationships {
    inventory_model_id = commercelayer_inventory_model.incentro_inventory_model.id
    stock_location_id  = commercelayer_stock_location.incentro_warehouse_location.id
  }
}