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
		fmt.Printf("Couldn't get repository: %v", err)
		return
	}
	//log.Print(repo)
	//tags, err := dapi.GetTags("", "nginx", 0, 0)
	//if err != nil {
	//	fmt.Printf("Couldn't get tags for repo: %v", err)
	//	return
	//}
	//log.Print(tags)
	//webhooks, err := dapi.GetWebhooks(dapi.GetUsername(), "nginx-proxy", 0,0)
	//if err != nil{
	//	fmt.Printf("Couldn't get webhooks for repo: %v", err)
	//}
	//log.Print(webhooks)
	//hook, err := dapi.CreateWebhook(dapi.GetUsername(), "nginx-proxy", "hook3", "https://google.com/")
	//if err != nil{
	//	fmt.Printf("Couldn't create webhook for repo: %v", err)
	//	return
	//}
	//log.Print(hook)
	//err = dapi.DeleteWebhook(dapi.GetUsername(), "nginx-proxy", "hook2")
	//if err != nil {
	//	fmt.Printf("Couldn't delete webhook for repo: %v", err)
	//	return
	//}
	//err = dapi.DeleteAllWebhooks(dapi.GetUsername(), "nginx-proxy")
	//if err != nil{
	//	fmt.Printf("Couldn't delete all webhooks: %v", err)
	//	return
	//}
	//err = dapi.SetRepositoryDescription(dapi.GetUsername(),"nginx-proxy","A fork of nginx-proxy", "#A multi-arch fork of nginx-proxy")
	//if err != nil {
	//	fmt.Printf("Couldn't set repository description: %v", err)
	//	return
	//}
	//dapi.DeleteAllWebhooks(dapi.GetUsername(), "nginx-proxy")
	//_, err = dapi.CreateWebhook(dapi.GetUsername(), "nginx-proxy", "basic2", "https://google.com")
	//if err != nil {
	//	fmt.Printf("Couldn't set webhook's url: %v", err)
	//	return
	//}
	//_, err = dapi.CreateOwnRepository("testing-repo", false, "Testing repo", "#Testing repo")
	//if err != nil{
	//	fmt.Printf("Could not create  repository: %v", err)
	//	return
	//}
	//err = dapi.DeleteOwnRepository("testing-repo")
	//if err != nil{
	//	fmt.Printf("Could not delete repository: %v", err)
	//	return
	//}

	repoLink := repo.GetGitRepo()
	log.Print(repoLink)

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
