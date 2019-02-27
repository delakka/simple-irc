package main

func main() {
	cfg := newConfig("./config.json")

	startServer(cfg)
}
