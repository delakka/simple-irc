package main

type Channel struct {
	Name  string
	Topic string
	Users map[*User]bool
}
