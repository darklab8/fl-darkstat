# cron job for:
# renew: docker run -v /var/lib/cerbot/:/var/www/certbot/ -v /var/lib/letsencrypt/:/etc/letsencrypt/ -it certbot/certbot:latest renew
resource "docker_image" "docker_cli" {
  name = "docker:27.3.1-cli"
}
locals {
  renewal_loop_seconds = 3600 * 24 * 3 # every 3 days
}
resource "docker_container" "renew" {
  name    = "nginx-renew"
  image   = docker_image.docker_cli.image_id
  restart = "always"
  tty     = true
  command = ["sh", "-c", "echo 'starting ${local.renewal_loop_seconds} cycle'; sleep ${local.renewal_loop_seconds}; docker run -v /var/lib/cerbot/:/var/www/certbot/ -v /var/lib/letsencrypt/:/etc/letsencrypt/ -it certbot/certbot:latest renew"]

  volumes {
    host_path      = "/var/run/docker.sock"
    container_path = "/var/run/docker.sock"
  }

  volumes {
    host_path      = "/var/lib/cerbot/"
    container_path = "/var/www/certbot/"
  }

  volumes {
    host_path      = "/var/lib/letsencrypt/"
    container_path = "/etc/letsencrypt/"
  }


  lifecycle {
    ignore_changes = [
      memory_swap,
      network_mode,
    ]
  }
}
