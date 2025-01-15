locals {
  envs = {
    FREELANCER_FOLDER = "/data/freelancer_folder" # Required

    RELAY_HOST        = var.RELAY_HOST # Required only for Discover Freelancer. Path to backend for PoBs data
    RELAY_LOOP_SECS    = "300" # Optional only for Discover Freelancer, how often to update PoB tab.

    DEV_ENV            = "true" # Optional: to optimize trade routes calcs to faster. optimize them better later
    DARKSTAT_LOG_LEVEL = "DEBUG" # Optional: for more debugging info
    UTILS_LOG_LEVEL    = "DEBUG" # Optional: for more debugging info

    SITE_ROOT = var.SITE_ROOT # Optional for when needing Site served from /fl-data-discovery/ route instead of just / (it is used to make Relay backend usable for Github pages frontend)
    FLDARKSTAT_HEADING = var.FLDARKSTAT_HEADING # Optional for phrases at the top of Darkstat interface
  }
}
