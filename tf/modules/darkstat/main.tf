resource "docker_image" "darkstat" {
  name         = "darkwind8/darkstat:${var.environment}"
  keep_locally = true
}

resource "docker_service" "darkstat" {
  name = "darkstat-${var.environment}"

  task_spec {
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
        memory_bytes = 1000 * 1000 * 6000 # 1 gb
      }
    }
  }
  endpoint_spec {
    mode = "vip"

    ports {
      target_port    = "8000"
      published_port = "8000"
    }

    ports {
      target_port    = "8080"
      published_port = "8080"
    }
  }
}
