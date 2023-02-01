resource "commercelayer_adyen_gateway" "incentro_adyen_gateway" {
  attributes {
    name                   = "Incentro Adyen Gateway"
    merchant_account       = "xxxx-yyyy-zzzz"
    api_key       		   = "xxxx-yyyy-zzzz"
    live_url_prefix        = "1797a841fbb37ca7-AdyenDemo"
  }
}