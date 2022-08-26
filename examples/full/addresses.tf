resource "commercelayer_address" "incentro-address" {
  attributes {
    business     = true
    company      = "Incentro"
    line_1       = "Van Nelleweg 1"
    zip_code     = "3044 BC"
    country_code = "NL"
    city         = "Rotterdam"
    phone        = "+31(0)10 20 20 544"
    state_code   = "ZH"
    metadata = {
      foo: "bar"
    }
  }

  relationships {

  }
}