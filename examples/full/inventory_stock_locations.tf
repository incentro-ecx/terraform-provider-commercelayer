resource "commercelayer_inventory_stock_location" "incentro_warehouse_location" {
  attributes {
    priority = 2
  }

  relationships {
    inventory_model_id = commercelayer_inventory_model.incentro_inventory_model.id
    stock_location_id  = commercelayer_stock_location.incentro_warehouse_location.id
  }
}

resource "commercelayer_inventory_stock_location" "incentro_backorder_location" {
  attributes {
    priority = 1
    on_hold  = true
  }

  relationships {
    inventory_model_id = commercelayer_inventory_model.incentro_inventory_model.id
    stock_location_id  = commercelayer_stock_location.incentro_backorder_location.id
  }
}