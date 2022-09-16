resource "commercelayer_merchant" "incentro_merchant" {
  attributes {
    name     = "Incentro Merchant"
    metadata = {
      foo : "bar"
    }
  }


  relationships {
    address = commercelayer_address.incentro_address.id
  }
}