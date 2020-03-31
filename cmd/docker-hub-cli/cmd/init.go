package cmd

import (
	"github.com/spf13/cobra"
)

const name = "docker-hub-cli"

var rootCmd = newRootCmd(name, "Helps you navigate the docker hub through the console.")

func newRootCmd(name, desc string) *cobra.Command {
	var c = &cobra.Command{Use: name}
	return c
}

func Execute() {
	_ = rootCmd.Execute()
}
