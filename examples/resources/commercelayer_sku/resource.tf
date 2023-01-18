resource "commercelayer_shipping_category" "incentro_shipping_category" {
  attributes {
    name = "Incentro Shipping Category"
  }
}

resource "commercelayer_sku" "incentro_sku" {
  attributes {
    name = "Incentro SKU"
    code = "TSHIRTMM000000FFFFFFXLXX"
  }
  relationships {
    shipping_category_id = commercelayer_shipping_category.incentro_shipping_category.id
  }
}