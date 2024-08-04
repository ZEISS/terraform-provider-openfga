#!/bin/bash

# Define the code to be written to the file
code="provider_installation {

  dev_overrides {
      \"registry.terraform.io/zeiss/terraform-provider-openfga\" = \"$(go env GOBIN)\"
  }

  # For all other providers, install them directly from their origin provider
  # registries as normal. If you omit this, Terraform will _only_ use
  # the dev_overrides block, and so no other providers will be available.
  direct {}
}"

# Write the code to the ~/.terraformrc file
echo "$code" > ~/.terraformrc

go build main.go
