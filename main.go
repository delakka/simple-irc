package main

import (
	"bufio"
	"encoding/json"
	"log"
	"net"
	"os"
	"strings"
)

func manageUser(user *User) {
	reader := bufio.NewReader(user.Conn)
	for {
		message, _ := reader.ReadString('\n')
		if len(message) == 0 {
			continue
		}
		message = strings.TrimSpace(message)
		log.Print("***MSG: ", string(message))
	}
}

func main() {
	file, err := os.Open("./config.json")
	defer file.Close()
	check(err)

	decoder := json.NewDecoder(file)
	conf := Config{}
	err = decoder.Decode(&conf)
	check(err)

	listener, err := net.Listen("tcp", conf.Port)
	defer listener.Close()
	check(err)
	log.Print("Listening on port ", conf.Port)

	connCh := make(chan net.Conn)
	go func() {
		for {
			conn, err := listener.Accept()
			check(err)
			log.Print("A new user connected.")
			connCh <- conn
		}
	}()

	for {
		conn := <-connCh
		user := &User{Conn: conn, Channels: make(map[*Channel]bool)}
		go manageUser(user)
	}
}
