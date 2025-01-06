resource "docker_image" "docker_cli" {
  name = "docker:27.3.1-cli"
}
locals {
  restart_seconds = 3600 * 12 # every 12 hours
}
// ensures to catch new Discovery version
resource "docker_container" "cron_restart" {
  name    = "darkstat-cron-restart-${var.environment}"
  image   = docker_image.docker_cli.image_id
  restart = "always"
  tty     = true
  command = ["sh", "-c", "echo 'starting ${local.restart_seconds} cycle'; sleep ${local.restart_seconds}; docker service update --force darkstat-${var.environment}"]

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
