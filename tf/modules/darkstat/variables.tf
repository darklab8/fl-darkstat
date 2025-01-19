locals {
  envs = merge({
    FREELANCER_FOLDER = "/data/freelancer_folder" # Required

    DARKSTAT_LOG_LEVEL = "DEBUG" # Optional: for more debugging info
    UTILS_LOG_LEVEL    = "DEBUG" # Optional: for more debugging info

    SITE_ROOT          = var.SITE_ROOT # Optional for when needing Site served from /fl-data-discovery/ route instead of just / (it is used to make Relay backend usable for Github pages frontend)
    SITE_HOST          = "https://${var.stat_prefix}.${var.zone}"
    FLDARKSTAT_HEADING = var.FLDARKSTAT_HEADING # Optional for phrases at the top of Darkstat interface
    },
    var.RELAY_HOST != null ? {
      RELAY_HOST      = var.RELAY_HOST # Required only for Discover Freelancer which have frontend at Github Pages. Path to backend for PoBs data
      RELAY_LOOP_SECS = "300"          # Optional only for Discover Freelancer, how often to update PoB tab.
    } : {},
    var.password != null ? { DARKCORE_PASSWORD = var.password } : {},
    var.secret != null ? { DARKCORE_SECRET = var.secret } : {},
  )
}
