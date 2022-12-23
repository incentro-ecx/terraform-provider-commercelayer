resource "commercelayer_stripe_gateway" "incentro_stripe_gateway" {
  attributes {
    name  = "Incentro Stripe Gateway"
    login = "SecretPassword"
    metadata = {
      foo : "bar"
    }
  }
}