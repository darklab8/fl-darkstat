locals {
  alloy_config = file("${path.module}/epbf.cfg.alloy")
  alloy_envs = {
    OTEL_EXPORTER_OTLP_ENDPOINT = "http://alloy-traces:4318"
    OTEL_GO_AUTO_TARGET_EXE="/shared/main"
    OTEL_SERVICE_NAME="${var.environment}-darkstat-epbf"
    OTEL_PROPAGATORS="tracecontext,baggage"
  }
}

resource "docker_image" "epbf" {
  name         = "otel/autoinstrumentation-go:v0.21.0"
  keep_locally = true
}

resource "docker_volume" "epbf_data" {
  name = "${var.environment}-darkstat-epbf"
}

resource "docker_container" "epbf" {
  name    = "${var.environment}-darkstat-epbf"
  image   = docker_image.epbf.name
  env     = [for k, v in local.alloy_envs : "${k}=${v}"]
    networks_advanced {
      name    = data.docker_network.grafana.id
      aliases = ["${var.environment}-darkstat-epbf"]
    }
  restart = "always"
  privileged = true
  pid_mode = "host"

  log_opts = {
    "mode" : "non-blocking"
    "max-buffer-size" : "500m"
  }
#   entrypoint = ["sh", "-c"]
#   command = [join(" && ", [
#     "echo '${local.alloy_config}' > /etc/alloy/config.alloy",
#     "/bin/alloy run /etc/alloy/config.alloy --storage.path=/var/lib/alloy/data",
#   ])]

  mounts {
    target    = "/shared"
    source    = docker_volume.epbf_data.name
    type      = "volume"
    read_only = false
  }
  volumes {
    container_path    = "/host/proc"
    host_path    = "/proc"
    read_only = true
  }

  memory = 1000 # MBs
  lifecycle {
    ignore_changes = [
      memory_swap,
      network_mode,
      image,
    ]
  }
}
