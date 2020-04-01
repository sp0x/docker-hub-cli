package main

type Configuration struct {
	Auth AuthConfiguration
}

type AuthConfiguration struct {
	Username string
	Token    string
}
