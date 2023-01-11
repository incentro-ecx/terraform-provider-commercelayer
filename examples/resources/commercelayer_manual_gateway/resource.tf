resource "commercelayer_manual_gateway" "incentro_manual_gateway" {
  attributes {
    name = "Incentro Manual Gateway"
    metadata = {
      foo : "bar"
    }
  }
}