terraform {
  required_providers {
    grafana = {
      source = "grafana/grafana"
    }
    curl = {
      source = "anschoewe/curl"
    }
  }
}

data "external" "secrets" {
  program = ["pass", "personal/terraform/grafana"]
}

locals {
  grafana_password = data.external.secrets.result["grafana_password"]
  grafana_creds    = "admin:${local.grafana_password}"
  discord_webhook_url = data.external.secrets.result["discord_webhook_url"]
}

provider "grafana" {
  url  = "https://grafana.dd84ai.com/"
  auth = local.grafana_creds
}
