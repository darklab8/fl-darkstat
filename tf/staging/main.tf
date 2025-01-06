module "discovery" {
  source      = "../modules/discovery"
  environment = "staging"
}

module "darkstat" {
  source             = "../modules/darkstat"
  environment        = "staging"
  discovery_path     = module.discovery.discovery_path
  ipv4_address       = module.data_cluster.node_darklab.ipv4_address
  RELAY_HOST         = "https://darkrelay-staging.dd84ai.com"
  SITE_ROOT          = "/fl-darkstat/"
  FLDARKSTAT_HEADING = <<-EOT
  <span style="font-weight:1000;">DEV ENV</span> <a href="https://github.com/darklab8/fl-darkstat">fl-darkstat</a> for <a href="https://github.com/darklab8/fl-data-discovery">Freelancer Discovery</a>
  EOT
  darkstat_port      = 8001
  relay_port         = 8081
}

module "dns" {
  source = "../../../infra/tf/modules/cloudflare_dns"
  zone   = "dd84ai.com"
  dns_records = [{
    type    = "A"
    value   = module.data_cluster.node_darklab.ipv4_address
    name    = "darkstat-staging"
    proxied = false
    }, {
    type    = "A"
    value   = module.data_cluster.node_darklab.ipv4_address
    name    = "darkrelay-staging"
    proxied = false
    }
  ]
}
