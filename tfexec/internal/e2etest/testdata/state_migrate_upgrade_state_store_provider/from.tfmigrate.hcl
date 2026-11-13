state_store_provider{
  tfcoremock = {
    source = "hashicorp/tfcoremock"
    version = "9.0.0"
  }
}

from {
  state_store "tfcoremock_fs" {
    provider "tfcoremock" {}

    workspace_dir = "old-location"
  }
}
