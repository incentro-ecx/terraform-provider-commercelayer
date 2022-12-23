resource "commercelayer_checkout_com_gateway" "incentro_checkout_com_gateway" {
  attributes {
    name                   = "Incentro CheckoutCom Gateway Changed"
    country_code              = "EU"
    api_key              = "xxxx-yyyy-zzzz"
    api_secret          = "xxxx-yyyy-zzzz"

    metadata = {
      bar: "foo"
      testName: "{{.testName}}"
    }
  }
}