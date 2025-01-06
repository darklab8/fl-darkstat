resource "docker_image" "nginx" {
  name = "nginx"
  build {
    context = path.module
  }

  triggers = {
    dir_sha1 = sha1(join("", [for f in ["default.conf", "Dockerfile", "run.sh"] : filesha1("${path.module}/${f}")]))
  }
}

# command to create cert from https://phoenixnap.com/kb/letsencrypt-docker
# docker run -v /var/lib/cerbot/:/var/www/certbot/ -v /var/lib/letsencrypt/:/etc/letsencrypt/ -it certbot/certbot:latest certonly --webroot --webroot-path /var/www/certbot/ --dry-run -d darkstat.dd84ai.com


resource "docker_container" "nginx" {
  name         = "nginx"
  image        = docker_image.nginx.image_id
  network_mode = "host" # lazy solution i know. Fix to proper networknig later.
  restart      = "always"

  # ports aren't needed in network mode
  #   ports {
  #     internal = "80"
  #     external = "80"
  #   }
  #   ports {
  #     internal = "443"
  #     external = "443"
  #   }

  volumes {
    host_path      = "/var/lib/cerbot/"
    container_path = "/var/www/certbot/"
  }

  volumes {
    host_path      = "/var/lib/letsencrypt/"
    container_path = "/var/lib/letsencrypt/"
  }


  lifecycle {
    ignore_changes = [
      memory_swap,
      network_mode,
    ]
  }
}

# cron job for:
# renew: docker run -v /var/lib/cerbot/:/var/www/certbot/ -v /var/lib/letsencrypt/:/etc/letsencrypt/ -it certbot/certbot:latest renew
resource "docker_image" "certbot" {
  name = "docker:27.3.1-cli"
}
locals {
  renewal_loop_seconds = 3600 * 24 * 3 # every 3 days
}
resource "docker_container" "renew" {
  name    = "nginx-renew"
  image   = docker_image.certbot.image_id
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
