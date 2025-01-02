resource "docker_image" "discovery" {
  name = "discovery"
  build {
    context = path.module
  }
}

locals {
  host_path = "/var/lib/darklab/discovery-${var.environment}"
}

resource "docker_container" "discovery" {
  name  = "discovery"
  image = docker_image.discovery.image_id

  volumes {
    host_path      = local.host_path
    container_path = "/code"
  }

  restart = "always"

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
