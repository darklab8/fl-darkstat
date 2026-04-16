resource "docker_image" "disco_api" {
  name = "disco-api"
  build {
    context = path.module
  }

  triggers = {
    dir_sha1 = sha1(join("",
      [for f in ["Dockerfile", "main.go", "go.mod", "go.sum"] : filesha1("${path.module}/${f}")]
    ))
  }
}

data "docker_network" "caddy" {
  name = "caddy"
}

resource "docker_container" "disco_api" {
  name    = "disco-api"
  image   = docker_image.disco_api.image_id
  restart = "always"
  tty     = true

  networks_advanced {
    name    = data.docker_network.caddy.id
    aliases = ["disco-api"]
  }

  log_opts = {
    "max-file" : "3"
    "max-size" : "10m"
  }
  volumes {
    host_path      = "/tmp/disco_api_data"
    container_path = "/data"
  }

  dynamic "labels" {
    for_each = merge({
      "caddy_0"               = "disco-api.dd84ai.com"
      "caddy_0.reverse_proxy" = "{{upstreams 8000}}"
      },
    )
    content {
      label = labels.key
      value = labels.value
    }
  }

  lifecycle {
    ignore_changes = [
      memory_swap,
      network_mode,
    ]
  }
}

variable "ipv4_address" {
  type = string
}

module "dns" {
  source = "../../../../infra/tf/modules/cloudflare_dns"
  zone   = "dd84ai.com"
  dns_records = concat([{
    type  = "A"
    value = var.ipv4_address
    name  = "disco-api"
    }
    ],
  )
}
