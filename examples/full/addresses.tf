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
    metadata = {
      foo: "bar"
    }
  }

  relationships {

  }
}

resource "commercelayer_address" "incentro-address2" {
  attributes {
    business     = true
    company      = "Incentro2"
    line_1       = "Van Nelleweg 12"
    zip_code     = "3044 BC"
    country_code = "NL"
    city         = "Rotterdam"
    phone        = "+31(0)10 20 20 542"
    state_code   = "ZH"
    metadata = {
      foo: "bar"
    }
  }

  relationships {

  }
}

resource "commercelayer_address" "incentro-address3" {
  attributes {
    business     = true
    company      = "Incentro3"
    line_1       = "Van Nelleweg 13"
    zip_code     = "3044 BC"
    country_code = "NL"
    city         = "Rotterdam"
    phone        = "+31(0)10 20 20 543"
    state_code   = "ZH"
    metadata = {
      foo: "bar"
    }
  }

  relationships {

  }
}

resource "commercelayer_address" "example-address" {
  attributes {
      business = false
      first_name = "John"
      last_name = "Smith"
      company ="The Red Brand Inc"
      line_1 =  "2883 Geraldine Lane"
      city = "New York"
      zip_code = "10013"
      state_code = "NY"
      country_code = "US"
      phone = "(212) 646-338-1228"
    metadata = {
      foo: "bar"
    }
  }

  relationships {

  }
}

