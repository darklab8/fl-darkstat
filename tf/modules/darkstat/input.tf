variable "environment" {
  type = string
}

variable "tag" {
  type    = string
  default = null
}

variable "discovery_path" {
  type = string
}

variable "ipv4_address" {
  type = string
}

variable "RELAY_HOST" {
  type    = string
  default = null
}
variable "SITE_ROOT" {
  type = string
}
variable "FLDARKSTAT_HEADING" {
  type = string
}
variable "stat_prefix" {
  type = string
}
variable "relay_prefix" {
  type    = string
  default = null
}
variable "rpc_prefix" {
  type    = string
  default = null
}
variable "pprof_prefix" {
  type    = string
  default = null
}
variable "apigateway_prefix" {
  type    = string
  default = null
}
variable "rpc_port" {
  type    = number
  default = null
}
variable "zone" {
  type = string
}

variable "password" {
  type    = string
  default = null
}
variable "secret" {
  type    = string
  default = null
}
variable "enable_restarts" {
  description = "good idea to turn on for mods that periodically load updates, like Discovery or FLSR. No need for Vanilla"
  type        = bool
}
variable "disco_oauth" {
  description = "https://github.com/darklab8/fl-darkstat/pull/106"
  type        = bool
  default     = false
}
variable "is_discovery" {
  type = bool
}