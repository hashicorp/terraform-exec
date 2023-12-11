variable "test" {
  type = string
}

resource "terraform_data" "test" {
  input = var.test
}

output "test" {
  value = terraform_data.test.output
}
