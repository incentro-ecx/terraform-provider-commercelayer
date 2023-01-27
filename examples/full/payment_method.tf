resource "commercelayer_payment_method" "incentro_payment_method" {
  attributes {
    payment_source_type   = "AdyenPayment"
    currency_code          = "EUR"
    price_amount_cents     = 0
  }

  relationships {
    payment_gateway_id = commercelayer_adyen_gateway.incentro_adyen_gateway.id
  }
}