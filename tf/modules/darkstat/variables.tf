locals {
  envs = merge({
    CONFIGS_FREELANCER_FOLDER = "/data/freelancer_folder" # Required

    DARKSTAT_LOG_LEVEL = "DEBUG" # Optional: for more debugging info
    UTILS_LOG_LEVEL    = "DEBUG" # Optional: for more debugging info

    DARKSTAT_SITE_ROOT          = var.SITE_ROOT # Optional for when needing Site served from /fl-data-discovery/ route instead of just / (it is used to make Relay backend usable for Github pages frontend)
    DARKSTAT_SITE_HOST          = "https://${var.stat_prefix}.${var.zone}"
    DARKSTAT_FLDARKSTAT_HEADING = var.FLDARKSTAT_HEADING # Optional for phrases at the top of Darkstat interface

    // grpc debugging
    GRPC_TRACE                  = "all"
    GRPC_VERBOSITY              = "DEBUG"
    GRPC_GO_LOG_SEVERITY_LEVEL  = "info"
    GRPC_GO_LOG_VERBOSITY_LEVEL = "6"
    },
    var.RELAY_HOST != null ? {
      DARKSTAT_RELAY_HOST      = var.RELAY_HOST # Required only for Discover Freelancer which have frontend at Github Pages. Path to backend for PoBs data
      DARKSTAT_RELAY_LOOP_SECS = "60"          # Optional only for Discover Freelancer, how often to update PoB tab.
    } : {},
    var.password != null ? {
      DARKCORE_PASSWORD = var.password
    } : {},
    var.secret != null ? {
      DARKCORE_SECRET = var.secret
    } : {},
    var.apigateway_prefix != null ? {
      DARKSTAT_GRPCGATEWAY_URL = "https://${var.apigateway_prefix}.${var.zone}/"
    } : {},
    var.disco_oauth == true ? {
      DARKCORE_DISCO_OAUTH = "true"
    } : {},
  )
}
