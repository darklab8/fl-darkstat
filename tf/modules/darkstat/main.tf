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

data "docker_network" "grafana" {
  name = "grafana"
}

resource "docker_service" "darkstat" {
  name = "${var.environment}-darkstat-app"

  task_spec {
    networks_advanced {
      name    = docker_network.network.id
      aliases = ["darkstat"]
    }
    networks_advanced {
      name    = data.docker_network.caddy.id
      aliases = ["${var.environment}-darkstat"]
    }
    networks_advanced {
      name    = data.docker_network.grafana.id
      aliases = ["${var.environment}-darkstat"]
    }
    log_driver {
      name = "json-file"

      options = {
        "max-file": "3"
        "max-size": "10m"
      }
    }
    container_spec {
      image = docker_image.darkstat.name
      env   = local.envs
      #   args = ["sleep", "infinity"]
      labels {
        label = "prometheus"
        value = "true"
      }
      dynamic "labels" {
        for_each = merge({
          "caddy_0"               = "${var.stat_prefix}.${var.zone}"
          "caddy_0.reverse_proxy" = "{{upstreams 8000}}"
          },
          var.relay_prefix != null ? {
            "caddy_1"               = "${var.relay_prefix}.${var.zone}"
            "caddy_1.reverse_proxy" = "{{upstreams 8080}}"
          } : {},
          var.rpc_prefix != null ? {
            "caddy_2"                                  = "${var.rpc_prefix}.${var.zone}:443",
            "caddy_2.reverse_proxy"                    = "{{upstreams h2c 50051}}"
            "caddy_3"                                  = "${var.rpc_prefix}.${var.zone}:80"
            "caddy_3.reverse_proxy.to"                 = "{{upstreams 50051}}"
            "caddy_3.reverse_proxy.transport"          = "http"
            "caddy_3.reverse_proxy.transport.versions" = "h1 h2c"
          } : {},
          var.apigateway_prefix != null ? {
            "caddy_4"               = "${var.apigateway_prefix}.${var.zone}"
            "caddy_4.reverse_proxy" = "{{upstreams 8081}}"
          } : {},
          var.pprof_prefix != null ? {
            "caddy_5"               = "${var.pprof_prefix}.${var.zone}"
            "caddy_5.reverse_proxy" = "{{upstreams 6060}}"
          } : {},
        )
        content {
          label = labels.key
          value = labels.value
        }
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
        memory_bytes = 1000 * 1000 * 5000 # 3 gb
      }
    }
  }
  lifecycle {
    ignore_changes = [
      task_spec[0].restart_policy[0].window,
      task_spec[0].container_spec[0].image,
      # task_spec[0].container_spec[0].env,
    ]
  }
  # with usage of docker networking, this is no longer necessary

  update_config {
    parallelism       = 1
    delay             = "60s"
    failure_action    = "pause"
    monitor           = "30s"
  }

  endpoint_spec {
    mode = "vip"

    dynamic "ports" {
      for_each = var.rpc_port != null ? ["rpc_port"] : []
      content {
        target_port    = "50051"
        published_port = tostring(var.rpc_port)
      }
    }
    # ports {
    #   target_port    = "50051"
    #   published_port = tostring(var.rpc_port)
    # }
  }
}
