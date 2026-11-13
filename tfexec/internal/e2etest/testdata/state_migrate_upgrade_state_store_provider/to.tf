terraform {
  required_providers {
    tfcoremock = {
      source = "hashicorp/tfcoremock"
      version = ">= 9.0.0"
    }
  }
  state_store "tfcoremock_fs" {
    provider "tfcoremock" {}

    workspace_dir = "new-location"
  }
}

resource "terraform_data" "test" {
  input = "foo"
}
