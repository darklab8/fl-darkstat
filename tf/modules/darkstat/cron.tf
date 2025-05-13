resource "docker_image" "docker_cli" {
  name = "docker:27.3.1-cli"
}
locals {
  restart_seconds = 3600 * 4 # every 4 hours
}
// ensures to catch new Discovery version
resource "docker_container" "cron_restart" {
  count   = var.enable_restarts ? 1 : 0
  name    = "${var.environment}-darkstat-cron-restart"
  image   = docker_image.docker_cli.image_id
  restart = "always"
  tty     = true
  command = ["sh", "-c", "echo 'starting ${local.restart_seconds} cycle'; sleep ${local.restart_seconds}; docker service update --force darkstat-${var.environment}"]
  log_opts = {
    "mode" : "non-blocking"
    "max-buffer-size" : "500m"
  }
  volumes {
    host_path      = "/var/run/docker.sock"
    container_path = "/var/run/docker.sock"
  }

  lifecycle {
    ignore_changes = [
      memory_swap,
      network_mode,
    ]
  }

}
