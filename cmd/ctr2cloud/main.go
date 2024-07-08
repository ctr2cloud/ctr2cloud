package main

import (
	"os"

	"github.com/ctr2cloud/ctr2cloud/cmd/ctr2cloud/raw"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(raw.Cmd)
}

var rootCmd = &cobra.Command{
	Use:          "ctr2cloud",
	Short:        "ctr2cloud is a cloud provider abstraction",
	SilenceUsage: true,
}

func main() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
