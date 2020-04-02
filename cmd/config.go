package main

import (
	"github.com/mitchellh/go-homedir"
	log "github.com/sirupsen/logrus"
	"github.com/sp0x/docker-hub-cli/api"
	"github.com/spf13/viper"
	"os"
)

var configFile string

//type Configuration struct {
//	Auth AuthConfiguration
//}

type AuthConfiguration struct {
	Username string
	Token    string
}

func (ac *AuthConfiguration) IsValid() bool {
	return ac.Username != "" && ac.Token != ""
}

func getUnauthorizedDockerApi() *api.DockerApi {
	return api.NewApi("", "")
}

func getAuthConfig() (*AuthConfiguration, error) {
	var authCfg AuthConfiguration
	err := viper.UnmarshalKey("auth", &authCfg)
	if err != nil {
		return nil, err
	}
	return &authCfg, nil
}

func getAvailableDockerApi() *api.DockerApi {
	authCfg, _ := getAuthConfig()
	if authCfg != nil && authCfg.IsValid() {
		var dapi = api.NewApi(authCfg.Username, authCfg.Token)
		return dapi
	} else {
		return getUnauthorizedDockerApi()
	}
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
				log.Warningf("error while writing default config file: %s\n %v\n", configFile, err)
			}
		} else {
			log.Warningf("error while reading config file: %s\n %v\n", configFile, err)
		}
	}

}
