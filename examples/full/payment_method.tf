resource "commercelayer_payment_method" "incentro_payment_method" {
  attributes {
    payment_source_type = "CreditCard"
    currency_code       = "EUR"
    price_amount_cents  = 1000
    metadata = {
      foo : "bar"
      testName : "{{.testName}}"
    }
  }

  relationships {
    market_id = commercelayer_market.incentro_market.id
  }
}

resource "commercelayer_inventory_model" "incentro_inventory_model" {
  attributes {
    name                   = "Incentro Inventory Model"
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


resource "commercelayer_merchant" "incentro_merchant" {
  attributes {
    name = "Incentro Merchant"
  }

  relationships {
    address_id = commercelayer_address.incentro_address.id
  }
}

resource "commercelayer_price_list" "incentro_price_list" {
  attributes {
    name          = "Incentro Price List"
    currency_code = "EUR"
  }
}

resource "commercelayer_market" "incentro_market" {
  attributes {
    name              = "Incentro Market"
    facebook_pixel_id = "pixel"
  }

  relationships {
    inventory_model_id = commercelayer_inventory_model.incentro_inventory_model.id
    merchant_id        = commercelayer_merchant.incentro_merchant.id
    price_list_id      = commercelayer_price_list.incentro_price_list.id
  }
}
