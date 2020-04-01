package main

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

const name = "docker-hub-cli"

var rootCmd = newRootCmd(name, "Helps you navigate the docker hub through the console.")

func newRootCmd(name, desc string) *cobra.Command {
	var c = &cobra.Command{Use: name, Short: desc}
	return c
}

func init() {
	//Ran after cobra is done.
	cobra.OnInitialize(initConfig)
	//We define our flags and configuration settings.
	//Cobra supports persistent flags which if defined here will be global for the whole app.
	rootCmd.PersistentFlags().StringVar(&configFile, "config", "", "config file (default is ~/.docker-hub-cli.yml")
	rootCmd.PersistentFlags().BoolP("verbose", "v", false, "Enable verbose mode")
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
