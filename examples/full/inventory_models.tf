resource "commercelayer_inventory_model" "incentro_inventory_model" {
  attributes {
    name                   = "Incentro Inventory Model"
    stock_locations_cutoff = 2
    strategy               = "split_shipments"
  }
}