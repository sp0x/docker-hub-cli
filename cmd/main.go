package main

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/sp0x/docker-hub-cli/api"
	"os"
)

func main() {
	dapi := api.NewApi()
	err := dapi.Login(os.Getenv("DOCKER_USER"), os.Getenv("DOCKER_PASS"))
	if err != nil {
		fmt.Println("Couldn't log in, try again.")
		return
	}
	//me, err := dapi.GetMyUser()
	//if err != nil {
	//	fmt.Printf("Couldn't get user: %v", err)
	//	return
	//}
	//log.Print(me)
	//repos, err := dapi.MyRepositories()
	//if err != nil {
	//	fmt.Printf("Couldn't get repositories")
	//	return
	//}
	//log.Print(repos)
	//repos, err := dapi.GetRepositoriesStarred(dapi.GetUsername(), 0, 0)
	//if err != nil {
	//	fmt.Printf("Couldn't get repositories")
	//	return
	//}
	//log.Print(repos)

	//repo, err := dapi.GetMyRepository("nginx-proxy")
	//if err != nil {
	//	fmt.Printf("Couldn't get repository")
	//	return
	//}
	//log.Print(repo)

	repo, err := dapi.GetRepository("", "nginx")
	if err != nil {
		fmt.Printf("Couldn't get repository %v", err)
		return
	}
	log.Print(repo)
	repoLinks := repo.GetGitRepoLinks()
	log.Print(repoLinks)

	buildSettings, err := dapi.GetBuildSettings(dapi.GetUsername(), "nginx-proxy")
	if err != nil {
		fmt.Printf("Couldn't get build settings")
		return
	}
	log.Print(buildSettings)

	err = dapi.Logout()
	if err != nil {
		fmt.Println("Couldn't logout")
	}

}
