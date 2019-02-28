package main

import (
	"net"
	"sync"
)

// Client is an accepted connection to the server
type Client struct {
	Nick          string
	UserName      string
	HostName      string
	RealName      string
	Conn          net.Conn
	Channels      map[string]*Channel
	Mutex         sync.RWMutex
	Authenticated bool
	Registered    bool
	//Messages      chan string
}

func newClient(conn net.Conn) *Client {
	c := &Client{
		Conn:          conn,
		Channels:      make(map[string]*Channel),
		Authenticated: false,
		Registered:    false,
	}

	return c
}

func (c *Client) send(message string) {
	//c.Messages <- message
	c.Conn.Write([]byte(message + "\r\n"))
}
