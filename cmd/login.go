package main

import (
	"bufio"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"golang.org/x/crypto/ssh/terminal"
	"os"
	"strings"
)

func init() {
	rootCmd.AddCommand(&cobra.Command{
		Use:   "login",
		Short: "Log into your docker hub account",
		Run: func(cmd *cobra.Command, args []string) {
			var authCfg AuthConfiguration
			err := viper.UnmarshalKey("auth", &authCfg)
			if err != nil {
				log.Warning("Could not unmarshal configuration.")
				return
			}
			duser := authCfg.Username
			var dpass string
			if authCfg.Token != "" {
				fmt.Printf("Already loggedin as %s", duser)
				os.Exit(0)
			}
			if duser == "" || authCfg.Token == "" {
				fmt.Print("Username: ")
				reader := bufio.NewReader(os.Stdin)
				duser, _ = reader.ReadString('\n')
				duser = strings.TrimSpace(duser)
				fmt.Print("Password: ")
				password, _ := terminal.ReadPassword(int(os.Stdin.Fd()))
				dpass = string(password)
			}
			if duser == "" || dpass == "" {
				//No username given
				fmt.Println("No authentication given.")
				os.Exit(1)
			}

			err = dockerApi.Login(duser, dpass)
			if err != nil {
				fmt.Println("Couldn't log in, try again.")
				return
			}
			authCfg.Username = duser
			authCfg.Token = dockerApi.GetToken()
			viper.Set("auth", authCfg)
			_ = viper.WriteConfig()
			fmt.Printf("Logged in.")
		},
	})
}
