resource "docker_image" "discovery_dev" {
  name = "discovery-dev-${var.environment}"
  build {
    context = path.module
  }

  triggers = {
    dir_sha1 = sha1(join("", [for f in ["Dockerfile", "main.go"] : filesha1("${path.module}/${f}")]))
  }
}

locals {
  host_path = "/var/lib/darklab/discovery-${var.environment}"
}

data "external" "disco_dev_webhook" {
  program = ["pass", "personal/terraform/darkstat/discovery_dev_branch_webhook"]
}


variable "docker_network_grafana_id" { type = string }

resource "docker_container" "discovery" {
  name  = "${var.environment}-darkstat-disco.dev.data"
  image = docker_image.discovery_dev.image_id

  volumes {
    host_path      = local.host_path
    container_path = "/code"
  }
  log_opts = {
    "max-file" : "3"
    "max-size" : "10m"
  }
  env = [
    "DISCO_DEV_WEBHOOK=${data.external.disco_dev_webhook.result["webhook_url"]}",
    "MY_SERVER_WEBHOOK=${data.external.disco_dev_webhook.result["my_server_url"]}"
  ]

  restart = "always"

  labels {
    label = "prometheus"
    value = "true"
  }
  networks_advanced {
    name    = var.docker_network_grafana_id
    aliases = ["${var.environment}-darkstat-dev.data"]
  }
  healthcheck {
    test         = ["CMD", "/install/main", "health"]
    interval     = "14s"
    timeout      = "20s"
    retries      = 6
    start_period = "2m"
  }

  volumes {
    host_path      = "/var/run/docker.sock"
    container_path = "/var/run/docker.sock"
  }

  lifecycle {
    ignore_changes = [
      memory_swap,
      network_mode,
    ]
  }
}

output "freelancer_path" {
  value = local.host_path
}
