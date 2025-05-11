
locals {
  grafana-dashboards = {
    darkstat_dashboard = {
      json = file("${path.module}/dashboards/darkstat.json")
    }
  }
}

resource "grafana_dashboard" "dashboard" {
  for_each    = local.grafana-dashboards
  config_json = each.value.json
}
