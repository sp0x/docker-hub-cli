package main

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

func init() {
	rootCmd.AddCommand(&cobra.Command{
		Use:   "repo",
		Short: "View, Create, Delete repositories",
		Run: func(cmd *cobra.Command, args []string) {
			dapi, err := getAuthorizedDockerApi()
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			repos, err := dapi.GetMyRepositories()
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
