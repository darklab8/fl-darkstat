variable "environment" {
  type = string
}

variable "tag" {
  type = string
  default = null
}

variable "discovery_path" {
  type = string
}

variable "ipv4_address" {
  type = string
}

variable "RELAY_HOST" {
  type = string
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
  type = string
}
variable "zone" {
  type = string
}

variable "password" {
  type = string
  default = null
}
variable "secret" {
  type = string
  default = null
}