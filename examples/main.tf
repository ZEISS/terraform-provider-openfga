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

resource "openfga_model" "demo" {
  spec = "{\"schema_version\":\"1.1\",\"type_definitions\":[{\"type\":\"user\"},{\"type\":\"document\",\"relations\":{\"reader\":{\"this\":{}},\"writer\":{\"this\":{}},\"owner\":{\"this\":{}}},\"metadata\":{\"relations\":{\"reader\":{\"directly_related_user_types\":[{\"type\":\"user\"}]},\"writer\":{\"directly_related_user_types\":[{\"type\":\"user\"}]},\"owner\":{\"directly_related_user_types\":[{\"type\":\"user\"}]}}}}]}"
  store = {
    id = openfga_store.demo.id
  }
}

resource "openfga_tuple" "demo" {
  user     = "user:demo"
  relation = "reader"
  document = "document:demo"

  store = {
    id    = openfga_store.demo.id
    model = openfga_model.demo.id
  }
}
