terraform {
  required_providers {
    tfcoremock = {
      source = "hashicorp/tfcoremock"
      version = "9.0.0" # Same as the source state
    }
  }
  state_store "tfcoremock_fs" {
    provider "tfcoremock" {}

    workspace_dir = "new-location" # Used to be "old-location"
  }
}

resource "terraform_data" "test" {
  input = "foo"
}
