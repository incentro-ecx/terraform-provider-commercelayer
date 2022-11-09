resource "commercelayer_market" "incentro_market" {
  attributes {
    name = "Incentro Market"
    //TODO: check why these are considered invalid
    #    checkout_url        = "http://example.url"
    #    external_prices_url = "http://example.url"
    facebook_pixel_id = "pixel"
  }

  relationships {
    inventory_model_id = commercelayer_inventory_model.incentro_inventory_model.id
    merchant_id        = commercelayer_merchant.incentro_merchant.id
    price_list_id      = commercelayer_price_list.incentro_price_list.id
    customer_group_id  = commercelayer_customer_group.incentro_customer_group.id
    tax_calculator_id  = commercelayer_external_tax_calculator.incentro_external_tax_calculator.id
  }
}