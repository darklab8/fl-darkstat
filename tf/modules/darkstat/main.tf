resource "docker_network" "network" {
  name       = "darkstat-${var.environment}"
  attachable = true
  driver     = "overlay"
}

locals {
  tag = var.tag != null ? var.tag : var.environment
}

resource "docker_image" "darkstat" {
  name         = "darkwind8/darkstat:${local.tag}"
  keep_locally = true
}

data "docker_network" "caddy" {
  name = "caddy"
}

resource "docker_service" "darkstat" {
  name = "${var.environment}-darkstat-app"

  task_spec {
    networks_advanced {
      name = docker_network.network.id
    }
    networks_advanced {
      name = data.docker_network.caddy.id
    }

    container_spec {
      image = docker_image.darkstat.name
      env   = local.envs
      #   args = ["sleep", "infinity"]

      labels {
        label = "caddy_0"
        value = "${var.stat_prefix}.${var.zone}"
      }
      labels {
        label = "caddy_0.reverse_proxy"
        value = "{{upstreams 8000}}"
      }
      labels {
        label = "caddy_1"
        value = "${var.relay_prefix}.${var.zone}"
      }
      labels {
        label = "caddy_1.reverse_proxy"
        value = "{{upstreams 8080}}"
      }

      mounts {
        target    = "/data"
        source    = var.discovery_path
        type      = "bind"
        read_only = false

        bind_options {
          propagation = "rprivate"
        }
      }
      mounts { // darkstat socks
        target    = "/tmp/darkstat"
        source    = "/tmp/darkstat-${var.environment}"
        type      = "bind"
        read_only = false
        bind_options {
          propagation = "rprivate"
        }
      }
    }
    restart_policy {
      condition = "any"
      delay     = "20s"
    }
    resources {
      limits {
        memory_bytes = 1000 * 1000 * 3000 # 1 gb
      }
    }
  }
  lifecycle {
    ignore_changes = [
      task_spec[0].restart_policy[0].window,
      task_spec[0].container_spec[0].env,
    ]
  }
  # with usage of docker networking, this is no longer necessary
  # endpoint_spec {
  #   mode = "vip"

  #   ports {
  #     target_port    = "8000"
  #     published_port = tostring(var.darkstat_port)
  #   }

  #   ports {
  #     target_port    = "8080"
  #     published_port = tostring(var.relay_port)
  #   }
  # }
}
