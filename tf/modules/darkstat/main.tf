resource "docker_network" "network" {
  name       = "darkstat-${var.environment}"
  attachable = true
  driver     = "overlay"
}

resource "docker_image" "darkstat" {
  name         = "darkwind8/darkstat:${var.environment}"
  keep_locally = true
}

resource "docker_service" "darkstat" {
  name = "darkstat-${var.environment}"

  task_spec {
    networks_advanced {
      name = docker_network.network.id
    }
    container_spec {
      image = docker_image.darkstat.name
      env   = local.envs
      #   args = ["sleep", "infinity"]

      mounts {
        target    = "/data"
        source    = var.discovery_path
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
