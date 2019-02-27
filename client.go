package main

import (
	"bufio"
	"log"
	"net"
	"strings"
)

// Client is an accepted connection to the server
type Client struct {
	Conn net.Conn
}

func newClient(conn net.Conn) *Client {
	return &Client{Conn: conn}
}

func (c *Client) handle() {
	reader := bufio.NewReader(c.Conn)
	for {
		message, _ := reader.ReadString('\n')
		if len(message) == 0 {
			continue
		}
		message = strings.TrimSpace(message)
		log.Print("***MSG: ", string(message))
	}
}
