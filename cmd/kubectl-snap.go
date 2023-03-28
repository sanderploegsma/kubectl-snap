package main

import (
	"os"

	"github.com/sanderploegsma/kubectl-snap/pkg/cmd"
	"github.com/spf13/pflag"

	"k8s.io/cli-runtime/pkg/genericclioptions"
)

// version can be automatically set by ldflags on build
var version = "dev"

func main() {
	flags := pflag.NewFlagSet("kubectl-snap", pflag.ExitOnError)
	pflag.CommandLine = flags

	root := cmd.NewSnapCmd(genericclioptions.IOStreams{In: os.Stdin, Out: os.Stdout, ErrOut: os.Stderr}, version)
	if err := root.Execute(); err != nil {
		os.Exit(1)
	}
}
