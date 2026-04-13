
module "vanilla" {
  count       = 0
  source      = "../modules/vanilla"
  environment = "production"
}

module "darkstat_vanilla" {
  count          = 0 // we save 356 mb of RAM by not running it
  source         = "../modules/darkstat"
  environment    = "vanilla"
  tag            = "production-arm64"
  discovery_path = module.vanilla[0].freelancer_path
  ipv4_address   = module.data_cluster.node_darklab.ipv4_address

  SITE_ROOT           = "/fl-data-vanilla/"
  FLDARKSTAT_HEADING  = <<-EOT
  <a href="https://github.com/darklab8/fl-darkstat">Darkstat</a> from <a href="https://darklab8.github.io/blog/pet_projects.html#Freelancercommunity">DarkTools</a> for Freelancer Vanilla
  EOT
  DARKSTAT_MAP_BY_URL = "https://darklab8.github.io/fl-data-vanilla/map.html"

  stat_prefix     = "darkstat-vanilla"
  zone            = "dd84ai.com"
  enable_restarts = false
  is_discovery    = false
  args            = ["--stat-deals-on", "web"]
}
