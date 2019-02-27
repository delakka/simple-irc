package main

func main() {
	cfg := NewConfig("./config.json")

	server := NewServer(cfg)
	server.Run()
}
