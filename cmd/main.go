package main

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"os"
)

func main() {
	api := NewApi()
	err := api.login(os.Getenv("DOCKER_USER"), os.Getenv("DOCKER_PASS"))
	if err != nil {
		fmt.Println("Couldn't log in, try again.")
		return
	}
	err = api.logout()
	if err != nil {
		fmt.Println("Couldn't logout")
	}

	t := api.getBuildSettings("sp0x", "nginx-proxy")
	log.Print(t)
}
