variable "create_duration" {
  type    = string
  default = "60s"
}

variable "destroy_duration" {
  type    = string
  default = null
}

resource "time_sleep" "sleep" {
  create_duration  = var.create_duration
  destroy_duration = var.destroy_duration
}
