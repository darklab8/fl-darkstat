locals {
  envs = {
    FREELANCER_FOLDER = "/data/freelancer_folder"
    RELAY_HOST        = "https://darkrelay.dd84ai.com"

    DEV_ENV            = "true" # to optimize trade routes. optimize them better later
    DARKSTAT_LOG_LEVEL = "DEBUG"
    UTILS_LOG_LEVEL    = "DEBUG"
    RELAY_LOOP_SECS    = "300"

    SITE_ROOT_ACCEPTORS = "/fl-data-discovery/,/fl-darkstat/"
    FLDARKSTAT_HEADING = <<-EOT
    <a href="https://github.com/darklab8/fl-darkstat">Darkstat</a> from <a href="https://darklab8.github.io/blog/pet_projects.html#Freelancercommunity">DarkTools</a> for <a href="https://github.com/darklab8/fl-data-discovery">Freelancer Discovery</a>
    EOT
  }
}
