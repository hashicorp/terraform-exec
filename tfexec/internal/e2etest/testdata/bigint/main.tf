terraform {
  required_providers {
    random = {
      version = "3.1.3"
    }
  }
}

resource "random_integer" "bigint" {
  max  = 7227701560655103598
  min  = 7227701560655103597
  seed = 12345
}
