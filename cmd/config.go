package main

import (
	"github.com/mitchellh/go-homedir"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"os"
)

var configFile string

type Configuration struct {
	Auth AuthConfiguration
}

type AuthConfiguration struct {
	Username string
	Token    string
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
	viper.SetDefault("docker.registry", "registry.hub.docker.com")
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
