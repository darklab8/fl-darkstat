
terraform {
  required_providers {
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

provider "cloudflare" {
  api_token = data.external.secrets_cloudflare.result["token"]
}

module "data_cluster" {
  source = "../../../infra/tf/production/output/deserializer"
}

provider "docker" {
  host     = "ssh://root@${module.data_cluster.node_darklab.ipv4_address}:22"
  ssh_opts = ["-o", "StrictHostKeyChecking=no", "-o", "UserKnownHostsFile=/dev/null", "-i", "~/.ssh/id_rsa.darklab"]
}
