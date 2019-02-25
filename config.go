package main

type Config struct {
	Server   string
	Port     string
	Channels []string

	Users []struct {
		Name     string
		Password string
		Info     string
	}
}
