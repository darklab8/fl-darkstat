locals {
  envs = merge({
    CONFIGS_FREELANCER_FOLDER = "/data/freelancer_folder" # Required

    TYPELOG_LOG_JSON = "true"

    DARKSTAT_SITE_ROOT          = var.SITE_ROOT # Optional for when needing Site served from /fl-data-discovery/ route instead of just / (it is used to make Relay backend usable for Github pages frontend)
    DARKSTAT_SITE_HOST          = "https://${var.stat_prefix}.${var.zone}"
    DARKSTAT_FLDARKSTAT_HEADING = var.FLDARKSTAT_HEADING # Optional for phrases at the top of Darkstat interface

    CACHE_CONTROL = "true"
    UTILS_ENVIRONMENT      = var.environment
    OTLP_HTTP_ON = var.environment == "staging" ? "true" : "false"
    OTEL_EXPORTER_OTLP_ENDPOINT = "http://alloy-traces:4318"
    OTEL_SERVICE_NAME = "${var.environment}-darkstat-app"
    OTEL_TRACES_SAMPLER = "parentbased_always_on"

    // grpc debugging
    GRPC_TRACE                  = "all"
    GRPC_VERBOSITY              = "DEBUG"
    GRPC_GO_LOG_SEVERITY_LEVEL  = "info"
    GRPC_GO_LOG_VERBOSITY_LEVEL = "6"
    },
    var.RELAY_HOST != null ? {
      DARKSTAT_RELAY_HOST      = var.RELAY_HOST # Required only for Discover Freelancer which have frontend at Github Pages. Path to backend for PoBs data
      DARKSTAT_RELAY_LOOP_SECS = "60"           # Optional only for Discover Freelancer, how often to update PoB tab.
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
    var.environment == "production" ? {
      OTEL_TRACES_SAMPLER_ARG = "0.1"
    } : {},
  )
}
