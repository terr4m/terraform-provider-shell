package main

import (
	"context"
	"flag"
	"log"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/terr4m/terraform-provider-shell/internal/provider"
)

var (
	// These will be set by GoReleaser.
	version string = "dev"
	commit  string = "none"
)

func main() {
	var debug bool

	flag.BoolVar(&debug, "debug", false, "set to true to run the provider with support for debuggers like delve")
	flag.Parse()

	opts := providerserver.ServeOpts{
		Address: "registry.terraform.io/terr4m/shell",
		Debug:   debug,
	}

	err := providerserver.Serve(context.Background(), provider.New(version, commit), opts)

	if err != nil {
		log.Fatal(err.Error())
	}
}
