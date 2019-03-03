package main

import (
	"log"
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
	Messages      chan string
}

func newClient(conn net.Conn) *Client {
	c := &Client{
		Conn:          conn,
		Channels:      make(map[string]*Channel),
		Authenticated: false,
		Registered:    false,
		Messages:      make(chan string),
	}

	return c
}

func (c *Client) sendLoop() {
	for message := range c.Messages {
		c.Conn.Write([]byte(message + "\r\n"))
		log.Print("[>>] ", string(message))
	}
}

func (c *Client) send(message string) {
	c.Messages <- message
}

func (c *Client) joinChannel(ch *Channel) {
	c.Channels[ch.Name] = ch
	log.Print("Added channel ", ch.Name, " to client ", c.Nick)
}

func (c *Client) leaveChannel(ch *Channel) {
	if _, ok := c.Channels[ch.Name]; ok {
		delete(c.Channels, ch.Name)
	}
}

func (c *Client) setNick(nick string) {
	for _, ch := range c.Channels {
		delete(ch.Clients, c.Nick)
		ch.Clients[nick] = c
	}
	c.Nick = nick
}

func (c *Client) in(clients []*Client) bool {
	for _, v := range clients {
		if v.Conn == c.Conn {
			return true
		}
	}
	return false
}

func (c *Client) getIndex(clients []*Client) (int, bool) {
	for i, v := range clients {
		if v.Conn == c.Conn {
			return i, true
		}
	}
	return 0, false
}

func findClient(nick string, clients []*Client) (*Client, bool) {
	for _, c := range clients {
		if c.Nick == nick {
			return c, true
		}
	}
	return &Client{}, false

}
