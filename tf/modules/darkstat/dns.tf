module "dns" {
  source = "../../../../infra/tf/modules/cloudflare_dns"
  zone   = var.zone
  dns_records = concat([{
    type    = "A"
    value   = var.ipv4_address
    name    = var.stat_prefix
    proxied = false
    }
  ],
  var.rpc_prefix != null ? [ {
    type    = "A"
    value   = var.ipv4_address
    name    = var.rpc_prefix
    proxied = false
    }] : [],
  var.relay_prefix != null ? [ {
    type    = "A"
    value   = var.ipv4_address
    name    = var.relay_prefix
    proxied = false
    }] : [],
  )
}
