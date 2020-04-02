package main

import (
	"errors"
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"strings"
)

var dockerfileTag string

func init() {
	dfCommand := &cobra.Command{
		Use:     "dockerfile [username/repo]",
		Short:   "Get the dockerfile for a repository",
		Aliases: []string{"df"},
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) > 1 {
				return errors.New("only one username accepted")
			} else if len(args) < 1 {
				return errors.New("repository missing")
			}
			return nil
		},
		Run: getDockerfileCommand,
	}
	dfCommand.Flags().StringVarP(&dockerfileTag, "tag", "t", "", "Use this if you want to find a specifically tagged dockerfile.")
	rootCmd.AddCommand(dfCommand)
}

func getDockerfileCommand(cmd *cobra.Command, args []string) {
	dapi := getAvailableDockerApi()
	name := args[0]
	parts := strings.Split(name, "/")
	if len(parts) == 1 {
		parts = append(parts, "")
		parts[0], parts[1] = "library", parts[0]
	}
	repo, err := dapi.GetRepository(parts[0], parts[1])
	if err != nil {
		fmt.Printf("Could not fetch repository %s: %s\n", name, err)
		os.Exit(1)
	}
	if dockerfileTag == "" {
		dockerfile, err := repo.GetDockerfile(dapi)
		if err != nil {
			fmt.Printf("Could not fetch dockerfile for %s: %s", name, err)
			os.Exit(1)
		}
		if dockerfile != "" {
			fmt.Print(dockerfile)
			os.Exit(0)
		} else {
			fmt.Println("Dockerfile is empty")
			os.Exit(1)
		}

	} else {

	}

}
