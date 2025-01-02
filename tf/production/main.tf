module "discovery" {
  source      = "../modules/discovery"
  environment = "production"
}

module "darkstat" {
  source         = "../modules/darkstat"
  environment    = "production"
  discovery_path = module.discovery.discovery_path
  ipv4_address   = module.data_cluster.node_darklab.ipv4_address
}

module "nginx" {
  source = "../modules/docker_nginx"
}

module "dns" {
  source = "../../../infra/tf/modules/cloudflare_dns"
  zone   = "dd84ai.com"
  dns_records = [{
    type    = "A"
    value   = module.data_cluster.node_darklab.ipv4_address
    name    = "darkstat"
    proxied = false
    }, {
    type    = "A"
    value   = module.data_cluster.node_darklab.ipv4_address
    name    = "darkrelay"
    proxied = false
    }, {
    type    = "A"
    value   = module.data_cluster.node_darklab.ipv4_address
    name    = "test"
    proxied = false
    }
  ]
}