routes = {
  "route1" = {
    address_prefix = "0.0.0.0/0",
    next_hop_type  = "Internet"
  },
  "route2" = {
    address_prefix         = "10.0.0.0/8",
    next_hop_type          = "VirtualAppliance",
    next_hop_in_ip_address = "10.10.10.0"
  }
}
logical_product_service = "routes"
