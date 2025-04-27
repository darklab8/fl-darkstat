module "dns" {
  source = "../../../../infra/tf/modules/cloudflare_dns"
  zone   = var.zone
  dns_records = concat([{
    type    = "A"
    value   = var.ipv4_address
    name    = var.stat_prefix
    proxied = true
    }
    ],
    var.rpc_prefix != null ? [{
      type  = "A"
      value = var.ipv4_address
      name  = var.rpc_prefix
    }] : [],
    var.pprof_prefix != null ? [{
      type  = "A"
      value = var.ipv4_address
      name  = var.pprof_prefix
    }] : [],
    var.relay_prefix != null ? [{
      type  = "A"
      value = var.ipv4_address
      name  = var.relay_prefix
    }] : [],
    var.apigateway_prefix != null ? [{
      type  = "A"
      value = var.ipv4_address
      name  = var.apigateway_prefix
    }] : [],
  )
}
