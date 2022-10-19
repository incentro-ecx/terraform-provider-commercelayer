resource "commercelayer_external_gateway" "incentro_external_gateway" {
  attributes {
    name          = "incentro_external_gateway"
    authorize_url = "https://example.com"
    capture_url = "https://foo.com"
    void_url = "https://foo.com"
    refund_url = "https://example.com"
    token_url = "https://example.com"
    metadata = {
      foo : "bar"
    }
  }
}