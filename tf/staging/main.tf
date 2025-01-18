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

  stat_prefix  = "darkstat-staging"
  relay_prefix = "darkrelay-staging"
  zone         = "dd84ai.com"
}
