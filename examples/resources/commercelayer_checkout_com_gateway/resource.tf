resource "commercelayer_checkout_com_gateway" "incentro_checkout_com_gateway" {
  attributes {
    name                   = "Incentro CheckoutCom Gateway Changed"
    secret_key             = "xxxx-yyyy-zzzz"
    public_key             = "xxxx-yyyy-zzzz"

    metadata = {
      bar: "foo"
      testName: "{{.testName}}"
    }
  }
}