package main

import (
	"os"

	"github.com/sanderploegsma/kubectl-snap/pkg/cmd"
	"github.com/spf13/pflag"
)

// version can be automatically set by ldflags on build
var version = "dev"

func init() {
	cmd.RootCmd.Version = version
}

func main() {
	flags := pflag.NewFlagSet("kubectl-snap", pflag.ExitOnError)
	pflag.CommandLine = flags

	if err := cmd.RootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
