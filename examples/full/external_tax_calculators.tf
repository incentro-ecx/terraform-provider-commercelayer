resource "commercelayer_external_tax_calculator" "incentro_external_tax_calculator" {
  attributes {
    name          = "incentro_external_tax_calculator"
    tax_calculator_url = "https://example.com"
    metadata = {
      foo : "bar"
    }
  }
}