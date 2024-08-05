package main

import (
	"github.com/zeiss/terraform-provider-openfga/cmd"
)

func main() {
	err := cmd.Root.Execute()
	if err != nil {
		panic(err)
	}
}
