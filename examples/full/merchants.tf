resource "commercelayer_merchant" "incentro_merchant" {
  attributes {
    name     = "Incentro Merchant"
    metadata = {
      foo : "bar"
    }
  }

  relationships {
    address_id = commercelayer_address.incentro_address.id
  }
}