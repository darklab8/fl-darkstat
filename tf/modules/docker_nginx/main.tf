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

# renew: docker run -v /var/lib/cerbot/:/var/www/certbot/ -v /var/lib/letsencrypt/:/etc/letsencrypt/ -it certbot/certbot:latest renew

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
