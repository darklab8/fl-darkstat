locals {
  envs = {
    FREELANCER_FOLDER = "/data/freelancer_folder"
    RELAY_HOST        = var.RELAY_HOST

    DEV_ENV            = "true" # to optimize trade routes. optimize them better later
    DARKSTAT_LOG_LEVEL = "DEBUG"
    UTILS_LOG_LEVEL    = "DEBUG"
    RELAY_LOOP_SECS    = "300"

    SITE_ROOT = var.SITE_ROOT
    # How the heck to fix that.
    # SITE_ROOT_ACCEPTORS = "/fl-data-discovery/,/fl-darkstat/"
    FLDARKSTAT_HEADING = var.FLDARKSTAT_HEADING
  }
}
