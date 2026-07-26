resource "docker_image" "cron_restart" {
  name = "${var.environment}-darkstat-patch.watcher"
  build {
    context = "${path.module}/patch.watcher"
  }

  triggers = {
    dir_sha1 = sha1(join("",
      [for f in ["Dockerfile", "main.go", "go.mod", "go.sum"] : filesha1("${path.module}/patch.watcher/${f}")]
    ))
  }
}

// ensures to catch new Discovery version
resource "docker_container" "cron_restart" {
  count   = var.enable_restarts ? 1 : 0
  name    = "${var.environment}-darkstat-patch.watcher"
  image   = docker_image.cron_restart.image_id
  restart = "always"
  tty     = true
  log_opts = {
    "max-file" : "3"
    "max-size" : "10m"
  }
  volumes {
    host_path      = "/var/run/docker.sock"
    container_path = "/var/run/docker.sock"
  }
  volumes {
    host_path      = var.discovery_path
    container_path = "/data"
  }
  labels {
    label = "prometheus"
    value = "true"
  }
  networks_advanced {
    name    = var.docker_network_grafana_id
    aliases = ["${var.environment}-darkstat-patch.watcher"]
  }
  healthcheck {
    test         = ["CMD", "/code/main", "health"]
    interval     = "14s"
    timeout      = "20s"
    retries      = 6
    start_period = "2m"
  }
  env = [
    "ENVIRONMENT=${var.environment}",
    "UTILS_USERAGENT=darkwind/1.0",
    "DARKMAP_REFRESH_GH_TOKEN=${var.trigger_darkmap_refresh_key}",
  ]
  lifecycle {
    ignore_changes = [
      memory_swap,
      network_mode,
    ]
  }
}
