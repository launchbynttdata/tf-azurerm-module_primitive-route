locals {
  route_table_name    = module.resource_names["route_table"].standard
  resource_group_name = module.resource_names["resource_group"].standard
  route_names         = { for k, v in var.routes : k => "${module.resource_names["route"].standard}${k}" }

  routes = {
    for k, v in var.routes : k => {
      name                   = local.route_names[k]
      resource_group_name    = module.resource_group.name
      route_table_name       = module.route_table.name
      address_prefix         = v.address_prefix
      next_hop_type          = v.next_hop_type
      next_hop_in_ip_address = v.next_hop_in_ip_address
    }
  }
}
