package main

import (
	"errors"
	"fmt"
	"github.com/sp0x/docker-hub-cli/api"
	"github.com/spf13/cobra"
	"os"
	"strings"
	"text/tabwriter"
)

func init() {
	rootCmd.AddCommand(&cobra.Command{
		Use:   "repo [username or username/repo]",
		Short: "View, Create, Delete repositories",
		Long:  "Use this to explore repositories or to manage them. If no username is given then the logged in user is used.",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) > 1 {
				return errors.New("only one username accepted")
			}
			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			var dapi *api.DockerApi
			var err error
			var repos []api.UserRepository
			isSpecificRepo := len(args) > 0 && strings.Contains(args[0], "/")
			if isSpecificRepo {
				showRepositoryDetails(args[0])
				return
			} else if len(args) > 0 {
				dapi = getAvailableDockerApi()
				repos, err = dapi.GetRepositories(args[0])
			} else {
				dapi = getAvailableDockerApi()
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

func showRepositoryDetails(fullName string) {
	dapi := getAvailableDockerApi()
	parts := strings.SplitN(fullName, "/", 2)
	repo, err := dapi.GetRepository(parts[0], parts[1])
	if err != nil {
		fmt.Printf("Could not fetch %s: %v", fullName, err)
	}
	gitRepo := repo.GetGitRepo()
	tags, err := dapi.GetTagsFromRepo(repo, 0, 0)
	//dockerfileContent, _ := dapi.GetDockerfileContents(parts[0], parts[1])
	//log.Print(dockerfileContent)
	fmt.Println(fullName)
	fmt.Println(repo.Description)
	fmt.Printf("Pulls: %d	Stars: %d\n", repo.PullCount, repo.StarCount)
	fmt.Printf("Last updated: %s\n", repo.LastUpdated)
	fmt.Printf("Git repo: %s\n", gitRepo)
	w := new(tabwriter.Writer)
	w.Init(os.Stdout, 0, 8, 0, '\t', 0)
	for _, tag := range tags {
		if repo.IsMarkdowned() {
			dockerfile, _ := repo.GetTaggedDockerfile(dapi, tag.Name, true)
			//dir, _ := repo.GetTaggedRepositoryDirectory(dapi, tag.Name, true)
			_, _ = fmt.Fprintf(w, "#%s\tBy: %s on %s\tDockerfile: %s\n", tag.Name, tag.LastUpdaterUsername, tag.LastUpdated, dockerfile)
		} else {
			_, _ = fmt.Fprintf(w, "Tag: %s\n", tag.Name)
		}
	}
	_ = w.Flush()
}
