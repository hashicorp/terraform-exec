terraform {
  backend "local" {
    path = "dst.tfstate"
  }
}

resource "terraform_data" "test" {
  input = "foo"
}
