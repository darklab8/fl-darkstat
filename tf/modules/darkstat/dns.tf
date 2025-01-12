module "dns" {
  source = "../../../../infra/tf/modules/cloudflare_dns"
  zone   = var.zone
  dns_records = [{
    type    = "A"
    value   = var.ipv4_address
    name    = var.stat_prefix
    proxied = false
    }, {
    type    = "A"
    value   = var.ipv4_address
    name    = var.relay_prefix
    proxied = false
    }
  ]
}
