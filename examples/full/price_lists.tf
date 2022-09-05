resource "commercelayer_price_list" "incentro_price_list" {
  attributes {
    name     = "incentro_group"
    currency_code = "EUR"
    metadata = {
      foo: "bar"
    }
  }
}