variable "in" {
  type    = string
  default = "default"
}

output "out" {
  value = var.in
}
