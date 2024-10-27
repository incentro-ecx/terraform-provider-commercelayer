resource "commercelayer_inventory_model" "incentro_inventory_model" {
  attributes {
    name                        = "Incentro Inventory Model"
    stock_locations_cutoff      = 2
    strategy                    = "split_shipments"
    manual_stock_decrement      = true
    stock_reservation_cutoff    = 4000
    put_stock_transfers_on_hold = true
  }
}