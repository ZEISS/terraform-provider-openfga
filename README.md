# Terraform Provider OpenFGA

[![Test & Build](https://github.com/zeiss/terraform-provider-openfga/actions/workflows/main.yml/badge.svg)](https://github.com/zeiss/terraform-provider-openfga/actions/workflows/main.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/zeiss/terraform-provider-openfga)](https://goreportcard.com/report/github.com/zeiss/terraform-provider-openfga)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Taylor Swift](https://img.shields.io/badge/secured%20by-taylor%20swift-brightgreen.svg)](https://twitter.com/SwiftOnSecurity)

:warning: **This provider is still in development and should not be used in production environments.**

## Example

```hcl
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
```

## Resources

- `openfga_store`
- `openfga_model`
- `openfga_tuple`

## Development

Run the following command to build the provider

```shell
sh ./scripts/setupDev.sh
```

Setup [Terraform](https://learn.hashicorp.com/tutorials/terraform/install-cli) and run the following command to test the provider

```shell
terraform init
```

Go to `./examples` and run the following command to test the provider

```bash
terraform plan
terraform apply
```

## License

[MIT](/LICENSE)