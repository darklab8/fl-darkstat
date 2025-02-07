resource "docker_image" "vanilla" {
  name = "vanilla"
  build {
    context = path.module
  }

  triggers = {
    dir_sha1 = sha1(join("", [for f in ["Dockerfile", "entrypoint.sh"] : filesha1("${path.module}/${f}")]))
  }
}

locals {
  host_path = "/var/lib/darklab/vanilla-${var.environment}"
}

resource "docker_container" "vanilla" {
  name  = "${var.environment}-vanilla-data"
  image = docker_image.vanilla.image_id

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

output "vanilla_path" {
  value = local.host_path
}
