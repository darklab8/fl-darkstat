
terraform {
  required_providers {
    hcloud = {
      source  = "hetznercloud/hcloud"
      version = ">= 1.35.2"
    }
    cloudflare = {
      source  = "cloudflare/cloudflare"
      version = ">=3.7.0"
    }
    docker = {
      source  = "kreuzwerker/docker"
      version = ">=3.0.2"
    }
  }
}

data "external" "secrets_cloudflare" {
  program = ["pass", "personal/terraform/cloudflare/dd84ai"]
}

data "external" "secrets_hetzner" {
  program = ["pass", "personal/terraform/hetzner/production"]
}

provider "hcloud" {
  token = data.external.secrets_hetzner.result["token"]
}

provider "cloudflare" {
  api_token = data.external.secrets_cloudflare.result["token"]
}

module "server" {
  source = "../../../infra/tf/modules/hetzner_server/data"
  name   = "node-arm"
}

provider "docker" {
  host     = "ssh://root@${module.server.ipv4_address}:22"
  ssh_opts = ["-o", "StrictHostKeyChecking=no", "-o", "UserKnownHostsFile=/dev/null", "-i", "~/.ssh/id_rsa.darklab"]
}
