package main

import (
	"log"
	"net"
)

func startServer(cfg *Config) {
	ln, err := net.Listen("tcp", cfg.Port)
	defer ln.Close()
	check(err)
	log.Print("Listening on port ", cfg.Port)

	acceptLoop(ln)

	// for {
	// 	conn := <-connCh
	// 	user := &User{Conn: conn, Channels: make(map[string]*Channel)}
	// 	go manageUser(user)
	// }
}

func acceptLoop(ln net.Listener) {
	// connCh := make(chan net.Conn)
	for {
		conn, err := ln.Accept()
		check(err)
		log.Print("A new user connected.")
		// connCh <- conn
		client := newClient(conn)
		go client.handle()
	}
}
