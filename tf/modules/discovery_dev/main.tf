resource "docker_image" "discovery_dev" {
  name = "discovery-dev-${var.environment}"
  build {
    context = path.module
  }

  triggers = {
    dir_sha1 = sha1(join("", [for f in ["Dockerfile", "entrypoint.py"] : filesha1("${path.module}/${f}")]))
  }
}

locals {
  host_path         = "/var/lib/darklab/discovery-${var.environment}"
  disco_dev_webhook = data.external.disco_dev_webhook.result["webhook_url"]
}

data "external" "disco_dev_webhook" {
  program = ["pass", "personal/terraform/darkstat/discovery_dev_branch_webhook"]
}

resource "docker_container" "discovery" {
  name  = "${var.environment}-discovery-data"
  image = docker_image.discovery_dev.image_id

  volumes {
    host_path      = local.host_path
    container_path = "/code"
  }
  log_opts = {
    "max-file": "3"
    "max-size": "10m"
  }
  env = [
    "DISCO_DEV_WEBHOOK=${local.disco_dev_webhook}"
  ]

  restart = "always"

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
