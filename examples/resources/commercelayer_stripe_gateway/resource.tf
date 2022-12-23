resource "commercelayer_stripe_gateway" "incentro_stripe_gateway" {
  attributes {
    name  = "Incentro Stripe Gateway"
    login = "sk_live_xxxx-yyyy-zzzz"
    metadata = {
      foo : "bar"
    }
  }
}