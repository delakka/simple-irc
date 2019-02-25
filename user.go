package main

import "net"

type User struct {
	Nick     string
	Conn     net.Conn
	Channels map[*Channel]bool
}
