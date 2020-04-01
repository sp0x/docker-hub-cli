package main

import (
	"fmt"
	homedir "github.com/mitchellh/go-homedir"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
)

const name = "docker-hub-cli"

var configFile string

var rootCmd = newRootCmd(name, "Helps you navigate the docker hub through the console.")

func newRootCmd(name, desc string) *cobra.Command {
	var c = &cobra.Command{Use: name, Short: desc}
	return c
}

func init() {
	//
	cobra.OnInitialize(initConfig)
	//We define our flags and configuration settings.
	//Cobra supports persistent flags which if defined here will be global for the whole app.
	rootCmd.PersistentFlags().StringVar(&configFile, "config", "", "config file (default is ~/.docker-hub-cli.yml")
	rootCmd.PersistentFlags().BoolP("verbose", "v", false, "Enable verbose mode")
}

func initConfig() {
	if configFile != "" {
		viper.SetConfigFile(configFile)
	} else {
		//We load the default config file
		home, err := homedir.Dir()
		if err != nil {
			log.Errorf("Could not find home directory: %v", err)
			os.Exit(1)
		}
		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".docker-hub-cli")
	}
	//Read the environment
	viper.AutomaticEnv()
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			err = viper.SafeWriteConfig()
			if err != nil {
				log.Warning("error while writing default config file: %s\n %v\n", configFile, err)
			}
		} else {
			log.Warning("error while reading config file: %s\n %v\n", configFile, err)
		}
	}

}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
