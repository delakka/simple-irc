package main

import (
	"net"
	"sync"
)

// User is a client who connected to the server
type User struct {
	Nick     string
	Name     string
	Host     string
	RealName string
	Password string
	Conn     net.Conn
	Channels map[string]*Channel
	MsgQ     chan string
	Mutex    sync.RWMutex
}
