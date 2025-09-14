resource "docker_image" "docker_cli" {
  name = "docker:27.3.1-cli"
}
locals {
  restart_seconds = 3600 * 6 # every 6 hours
}
// ensures to catch new Discovery version
resource "docker_container" "cron_restart" {
  count   = var.enable_restarts ? 1 : 0
  name    = "${var.environment}-darkstat-cron-restart"
  image   = docker_image.docker_cli.image_id
  restart = "always"
  tty     = true
  command = ["sh", "-c", "echo 'starting ${local.restart_seconds} cycle'; sleep ${local.restart_seconds}; docker service update --force ${var.environment}-darkstat-app"]
  log_opts = {
    "max-file" : "3"
    "max-size" : "10m"
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
