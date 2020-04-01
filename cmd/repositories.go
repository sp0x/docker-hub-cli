package main

import (
	"errors"
	"fmt"
	"github.com/sp0x/docker-hub-cli/api"
	"github.com/spf13/cobra"
	"os"
	"strings"
)

func init() {
	rootCmd.AddCommand(&cobra.Command{
		Use:   "repo [username]",
		Short: "View, Create, Delete repositories",
		Long:  "Use this to explore repositories or to manage them. If no username is given then the logged in user is used.",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) > 1 {
				return errors.New("only one username accepted")
			}
			if len(args) > 0 && strings.Contains(args[0], "/") {
				return errors.New("username can't contain /")
			}
			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			var dapi *api.DockerApi
			var err error
			var repos []api.UserRepository
			if len(args) > 0 {
				dapi = newUnauthorizedDockerApi()
				repos, err = dapi.GetRepositories(args[0])
			} else {
				dapi, err = getAuthorizedDockerApi()
				if err != nil {
					fmt.Println(err)
					os.Exit(1)
				}
				repos, err = dapi.GetMyRepositories()
			}

			if err != nil {
				fmt.Printf("Error while listing repositories: %v", err)
				os.Exit(1)
			}
			for _, repo := range repos {
				fmt.Printf("%s/%s\n", repo.Namespace, repo.Name)
			}
		},
	})
}
