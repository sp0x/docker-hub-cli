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

var repoShowTags bool

func init() {
	reposCmd := &cobra.Command{
		Use:   "repo [username or username/repo]",
		Short: "View, Create, Delete repositories",
		Long:  "Use this to explore repositories or to manage them. If no username is given then the logged in user is used.",
		Args: func(cmd *cobra.Command, args []string) error {
			//if len(args) > 1 {
			//	return errors.New("only one username accepted")
			//}
			return nil
		},
		Run: reposCommand,
	}
	reposCmd.Flags().BoolVarP(&repoShowTags, "tags", "t", false, "Also shows all the tags in the repository")

	rmRepoCmd := &cobra.Command{
		Use:   "rm [repository]",
		Short: "Delete a repository",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) > 1 {
				return errors.New("only one repository accepted")
			} else if len(args) < 1 {
				return errors.New("repository is missing")
			}
			return nil
		},
		Run: rmRepoCommand,
	}
	createRepoCmd := &cobra.Command{
		Use:   "create [repository]",
		Short: "Delete a repository",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) > 1 {
				return errors.New("only one repository accepted")
			} else if len(args) < 1 {
				return errors.New("repository is missing")
			}
			return nil
		},
		Run: createRepoCommand,
	}
	reposCmd.AddCommand(rmRepoCmd)
	reposCmd.AddCommand(createRepoCmd)
	rootCmd.AddCommand(reposCmd)
}

func reposCommand(cmd *cobra.Command, args []string) {
	var dapi *api.DockerApi
	var err error
	var repos []api.UserRepository
	if len(args) > 0 {
		dapi = getAvailableDockerApi()
		for _, arg := range args {
			isSpecificRepo := strings.Contains(arg, "/")
			if isSpecificRepo {
				showRepositoryDetails(dapi, arg)
				fmt.Println("")
			} else if len(args) > 0 {
				repos, err = dapi.GetRepositories(arg)
			}
		}
	} else {
		dapi = getAvailableDockerApi()
		repos, err = dapi.GetMyRepositories()
		if err != nil {
			fmt.Printf("Error while listing repositories: %v\n", err)
			os.Exit(1)
		}
		for _, repo := range repos {
			fmt.Printf("%s/%s\n", repo.Namespace, repo.Name)
		}
	}

}

func rmRepoCommand(cmd *cobra.Command, args []string) {
	dapi := getAvailableDockerApi()
	if !dapi.IsAuthenticated() {
		fmt.Printf("You need to login first.\n")
		os.Exit(1)
	}
	repo := args[0]
	parts := strings.Split(repo, "/")
	//If no username is given then we'll use
	if len(parts) < 2 {
		parts = append(parts, "")
		parts[0], parts[1] = dapi.GetUsername(), parts[0]
	}
	err := dapi.DeleteRepository(parts[0], parts[1])
	if err != nil {
		fmt.Printf("Could not delete repository: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Removed repository %s/%s\n", parts[0], parts[1])
}

func createRepoCommand(cmd *cobra.Command, args []string) {
	dapi := getAvailableDockerApi()
	if !dapi.IsAuthenticated() {
		fmt.Printf("You need to login first.")
		os.Exit(1)
	}
	name := args[0]
	repo, err := dapi.CreateOwnRepository(name, false, "", "")
	if err != nil {
		fmt.Printf("\n")
		os.Exit(1)
	}
	fmt.Printf("Created repository: %s/%s", repo.Namespace, repo.Name)
}

func showRepositoryDetails(dapi *api.DockerApi, fullName string) {
	parts := strings.SplitN(fullName, "/", 2)
	repo, err := dapi.GetRepository(parts[0], parts[1])
	if err != nil {
		fmt.Printf("Could not fetch %s: %v", fullName, err)
	}
	gitRepo := repo.GetGitRepo()
	tags, err := dapi.GetTagsFromRepo(repo, 0, 0)

	fmt.Println(fullName)
	fmt.Println(repo.Description)
	fmt.Printf("Pulls: %d	Stars: %d\n", repo.PullCount, repo.StarCount)
	if repo.LastUpdated != nil {
		fmt.Printf("Last updated: %s\n", timeElapsedRightNow(*repo.LastUpdated, false))
	}
	if gitRepo != "" {
		fmt.Printf("Git repo: %s\n", gitRepo)
	}
	if repoShowTags {
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
}
