resource "docker_image" "discovery_dev" {
  name = "discovery-dev"
  build {
    context = path.module
  }

  triggers = {
    dir_sha1 = sha1(join("", [for f in ["Dockerfile", "entrypoint.sh"] : filesha1("${path.module}/${f}")]))
  }
}

locals {
  host_path = "/var/lib/darklab/discovery-${var.environment}"
}

resource "docker_container" "discovery" {
  name  = "${var.environment}-discovery-data"
  image = docker_image.discovery_dev.image_id

  volumes {
    host_path      = local.host_path
    container_path = "/code"
  }

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

output "discovery_path" {
  value = local.host_path
}
