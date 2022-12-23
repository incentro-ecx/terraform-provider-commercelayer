resource "commercelayer_inventory_model" "incentro_inventory_model" {
  attributes {
    name                   = "Incentro Inventory Model Return Location"
    stock_locations_cutoff = 2
    strategy               = "split_shipments"
  }
}

resource "commercelayer_address" "incentro_address" {
  attributes {
    business     = true
    company      = "Incentro"
    line_1       = "Van Nelleweg 1"
    zip_code     = "3044 BC"
    country_code = "NL"
    city         = "Rotterdam"
    phone        = "+31(0)10 20 20 544"
    state_code   = "ZH"
  }
}

resource "commercelayer_stock_location" "incentro_stock_location" {
  attributes {
    name         = "Incentro Warehouse Location"
    label_format = "PNG"
    suppress_etd = true
  }

  relationships {
    address_id = commercelayer_address.incentro_address.id
  }
}

resource "commercelayer_inventory_return_location" "incentro_return_location" {
  attributes {
    priority = 1
  }

  relationships {
    inventory_model_id = commercelayer_inventory_model.incentro_inventory_model.id
    stock_location_id  = commercelayer_stock_location.incentro_stock_location.id
  }
}