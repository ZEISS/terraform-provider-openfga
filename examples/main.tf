terraform {
  required_providers {
    openfga = {
      source = "zeiss/openfga"
    }
  }
}

provider "openfga" {
  api_url = "http://host.docker.internal:8080"
}

resource "openfga_store" "demo" {
  name = "demo"
}
