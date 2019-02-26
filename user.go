package main

import "net"

// An IRC user
type User struct {
	Nick     string
	Conn     net.Conn
	Channels map[*Channel]bool
}
