resource "commercelayer_delivery_lead_time" "incentro_delivery_lead_time" {
  attributes {
    min_hours = 10
    max_hours = 100
    metadata = {
      foo : "bar"
    }
  }

  relationships {
    stock_location_id = commercelayer_stock_location.incentro_warehouse_location.id
    shipping_method_id = commercelayer_shipping_method.incentro_shipping_method.id
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


resource "commercelayer_stock_location" "incentro_warehouse_location" {
  attributes {
    name         = "Incentro Warehouse Location"
    label_format = "PNG"
    suppress_etd = true
  }

  relationships {
    address_id = commercelayer_address.incentro_address.id
  }
}

resource "commercelayer_shipping_method" "incentro_shipping_method" {
  attributes {
    name                   = "Incentro Shipping Method"
    scheme                 = "flat"
    currency_code          = "EUR"
    price_amount_cents     = 1000
    free_over_amount_cents = 10000
    min_weight             = 0.50
    max_weight             = 10
    unit_of_weight         = "kg"
  }
}