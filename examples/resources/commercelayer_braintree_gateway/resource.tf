resource "commercelayer_braintree_gateway" "incentro_braintree_gateway" {
  attributes {
    name                   = "Incentro Braintree Gateway"
    merchant_account_id    = "xxxx-yyyy-zzzz"
    merchant_id            = "xxxx-yyyy-zzzz"
    public_key             = "xxxx-yyyy-zzzz"
    private_key            = "xxxx-yyyy-zzzz"

    metadata = {
      foo: "bar"
      testName: "{{.testName}}"
    }
  }
}