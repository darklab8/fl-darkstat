data "docker_network" "caddy" {
  name = "caddy"
}

data "docker_network" "grafana" {
  name = "grafana"
}

module "discovery" {
  source      = "../modules/disco.prod.data"
  environment = "production"
}

module "disco_api" {
  source                    = "../modules/disco.api"
  ipv4_address              = module.data_cluster.node_darklab.ipv4_address
  environment               = "production"
  docker_network_caddy_id   = data.docker_network.caddy.id
  docker_network_grafana_id = data.docker_network.grafana.id
}

module "darkstat" {
  source         = "../modules/darkstat"
  environment    = "production"
  tag            = "production-arm64"
  discovery_path = module.discovery.freelancer_path
  ipv4_address   = module.data_cluster.node_darklab.ipv4_address

  docker_network_caddy_id   = data.docker_network.caddy.id
  docker_network_grafana_id = data.docker_network.grafana.id

  SITE_ROOT           = "/fl-data-discovery/"
  FLDARKSTAT_HEADING  = <<-EOT
  <a href="https://github.com/darklab8/fl-darkstat">Darkstat</a> from <a href="https://darklab8.github.io/blog/pet_projects.html#Freelancercommunity">DarkTools</a> for <a href="https://github.com/darklab8/fl-data-discovery">Disco</a>
  EOT
  DARKSTAT_MAP_BY_URL = "https://darklab8.github.io/fl-data-discovery/map.html"

  stat_prefix                 = "darkstat"
  pprof_prefix                = "darkstat-pprof"
  zone                        = "dd84ai.com"
  is_discovery                = true
  is_discovery_production     = true
  enable_restarts             = true
  trigger_darkmap_refresh_key = local.trigger_darkmap_refresh_key

  replicas_count = 2
  extra_vars     = local.disco_extra_vars
  args           = ["--stat-deals-on", "web_cron"]
}

locals {
  disco_extra_vars = {
    CONFIGS_DISCO_BASES_FULL_URL = data.external.secrets_darkbot.result["SCRAPPY_BASE_URL"]
    DISABLE_DEV_MODE             = "true"
  }
}

resource "random_string" "random_password" {
  length  = 32
  special = false
}
resource "random_string" "random_secret" {
  length  = 32
  special = false
}

module "discovery_dev" {
  source                    = "../modules/disco.dev.data"
  environment               = "dev"
  docker_network_grafana_id = data.docker_network.grafana.id
}

module "darkstat_dev" {
  source         = "../modules/darkstat"
  environment    = "dev"
  tag            = "production-arm64"
  discovery_path = module.discovery_dev.freelancer_path
  ipv4_address   = module.data_cluster.node_darklab.ipv4_address

  docker_network_caddy_id   = data.docker_network.caddy.id
  docker_network_grafana_id = data.docker_network.grafana.id

  SITE_ROOT          = "/"
  FLDARKSTAT_HEADING = <<-EOT
  <a href="https://github.com/darklab8/fl-darkstat">Darkstat</a> from <a href="https://darklab8.github.io/blog/pet_projects.html#Freelancercommunity">DarkTools</a> for <a href="https://github.com/darklab8/fl-data-discovery">Disco</a>
  EOT

  stat_prefix = "darkstat-dev"
  zone        = "dd84ai.com"

  password                    = random_string.random_password.result
  secret                      = random_string.random_secret.result
  disco_oauth                 = true
  is_discovery                = true
  enable_restarts             = false
  trigger_darkmap_refresh_key = ""

  # extra_vars   = local.disco_extra_vars
  extra_vars = {
    "DARKCORE_LOG_LEVEL" = "DEBUG"
  }
  args = ["--stat-deals-on", "--map-on", "web"]
}
