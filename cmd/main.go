package main

import (
	log "github.com/sirupsen/logrus"
	"os"
)

func main() {
	api := NewApi()
	_ = api.login(os.Getenv("DOCKER_USER"), os.Getenv("DOCKER_PASS"))
	t := api.getBuildSettings("sp0x", "nginx-proxy")
	log.Print(t)
}
