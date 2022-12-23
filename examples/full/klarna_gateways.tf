resource "commercelayer_klarna_gateway" "incentro_klarna_gateway" {
  attributes {
    name                   = "Incentro Klarna Gateway Changed"
    country_code              = "EU"
    api_key              = "xxxx-yyyy-zzzz"
    api_secret          = "xxxx-yyyy-zzzz"

    metadata = {
      bar: "foo"
      testName: "{{.testName}}"
    }
  }
}