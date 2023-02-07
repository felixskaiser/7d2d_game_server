###############################################################################
# Required variables
###############################################################################

variable "project_id" {
  type = string
}

variable "default_region" {
  type = string
}

variable "default_zone" {
  type = string
}

variable "billing_account_id" {
  type = string
}

###############################################################################
# Optional variables
###############################################################################

variable "server_name" {
  type    = string
  default = "game-server-7d2d"
}

variable "machine_type" {
  type    = string
  default = "e2-medium"
}

variable "disk_size" {
  type    = string
  default = "100"
}