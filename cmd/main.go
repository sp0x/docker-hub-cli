package main

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/sp0x/docker-hub-cli/api"
	"os"
)

func main() {
	dapi := api.NewApi()
	err := dapi.login(os.Getenv("DOCKER_USER"), os.Getenv("DOCKER_PASS"))
	if err != nil {
		fmt.Println("Couldn't log in, try again.")
		return
	}
	me, err := dapi.getMyUser()
	if err != nil {
		fmt.Printf("Couldn't get user: %v", err)
		return
	}
	log.Print(me)
	err = dapi.logout()
	if err != nil {
		fmt.Println("Couldn't logout")
	}
	t := dapi.getBuildSettings("sp0x", "nginx-proxy")
	log.Print(t)
}
