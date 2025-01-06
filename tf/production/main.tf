module "discovery" {
  source      = "../modules/discovery"
  environment = "production"
}

module "darkstat" {
  source         = "../modules/darkstat"
  environment    = "production"
  discovery_path = module.discovery.discovery_path
  ipv4_address   = module.data_cluster.node_darklab.ipv4_address

  RELAY_HOST         = "https://darkrelay.dd84ai.com"
  SITE_ROOT          = "/fl-data-discovery/"
  FLDARKSTAT_HEADING = <<-EOT
  <a href="https://github.com/darklab8/fl-darkstat">Darkstat</a> from <a href="https://darklab8.github.io/blog/pet_projects.html#Freelancercommunity">DarkTools</a> for <a href="https://github.com/darklab8/fl-data-discovery">Freelancer Discovery</a>
  EOT
  darkstat_port      = 8000
  relay_port         = 8080
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
