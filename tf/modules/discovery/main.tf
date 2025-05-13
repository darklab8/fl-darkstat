resource "docker_image" "discovery" {
  name = "discovery-${var.environment}"
  build {
    context = path.module
  }
}

locals {
  host_path = "/var/lib/darklab/discovery-${var.environment}"
}

resource "docker_container" "discovery" {
  name  = "${var.environment}-discovery-data"
  image = docker_image.discovery.image_id

  volumes {
    host_path      = local.host_path
    container_path = "/code"
  }
  log_opts = {
    "mode" : "non-blocking"
    "max-buffer-size" : "500m"
  }
  restart = "always"

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
