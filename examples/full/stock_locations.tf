resource "commercelayer_stock_location" "incentro_stock_location" {
  attributes {
    name         = "Incentro Stock Location"
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