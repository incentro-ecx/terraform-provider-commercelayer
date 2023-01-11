resource "commercelayer_adyen_gateway" "incentro_adyen_gateway" {
  attributes {
    name = "Incentro Adyen Gateway"
    metadata = {
      foo : "bar"
    }
  }
}