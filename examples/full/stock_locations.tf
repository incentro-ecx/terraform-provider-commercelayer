resource "commercelayer_stock_location" "incentro_warehouse_location" {
  attributes {
    name         = "Incentro Warehouse Location"
    label_format = "PNG"
    suppress_etd = true
    metadata = {
      foo : "bar"
    }
  }

  relationships {
    address_id = commercelayer_address.incentro_address.id
  }
}

resource "commercelayer_stock_location" "incentro_backorder_location" {
  attributes {
    name         = "Incentro Backorder Location"
    label_format = "PNG"
    suppress_etd = true
    metadata = {
      foo : "bar"
    }
  }

  relationships {
    address_id = commercelayer_address.incentro_address.id
  }
}