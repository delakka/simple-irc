package main

import (
	"bufio"
	"log"
	"net"
	"strings"
)

// Server represents an IRC server
type Server struct {
	Listener net.Listener
	Channels map[string]*Channel
	Clients  []*Client
	Password string
	Port     string
	Host     string
}

// NewServer instantiates a server
func NewServer(cfg *Config) *Server {
	return &Server{
		Channels: make(map[string]*Channel),
		Clients:  make([]*Client, 0),
		Password: cfg.Password,
		Port:     cfg.Port,
		Host:     cfg.Server,
	}
}

// Run starts the server
func (s *Server) Run() {
	ln, err := net.Listen("tcp", s.Port)
	defer ln.Close()
	check(err)
	log.Print("Listening on port ", s.Port)
	s.Listener = ln

	go s.receiveLoop()
	s.acceptLoop()
}

func (s *Server) acceptLoop() {
	for {
		conn, err := s.Listener.Accept()
		check(err)
		log.Print("A new user connected with the remote IP: ", conn.RemoteAddr())
		client := newClient(conn)
		go client.sendLoop()
		s.addClient(client)
	}
}

func (s *Server) receiveLoop() {
	for {
		for _, c := range s.Clients {
			// buf := make([]byte, 1024)
			reader := bufio.NewReader(c.Conn)
			for {
				// log.Print("C: ", c.Conn.RemoteAddr().String())
				message, _ := reader.ReadString('\n')
				// log.Print("ML: ", len(message), " M: ", message)
				if len(message) == 0 {
					continue
				}
				message = strings.TrimSpace(message)
				log.Print("[<<] ", string(message))

				cmd, params := parseMessage(message)
				if _, ok := commands[cmd]; ok {
					log.Print(cmd, " command received from client ", c.Nick, " with params: ", strings.Join(params, " "))
					commands[cmd](params, c, s)
				}
			}
		}
	}
}

func (s *Server) isNickAvailable(nick string) bool {
	for _, v := range s.Clients {
		if v.Nick == nick {
			return false
		}
	}
	return true
}

func (s *Server) auth(password string) bool {
	if s.Password == password {
		return true
	}
	return false
}

func (s *Server) addClient(c *Client) {
	if c.in(s.Clients) {
		log.Print("Client was already added to the server")
		return
	}
	s.Clients = append(s.Clients, c)
	log.Print("Added client ", c.Nick, " to the server")
}

func (s *Server) removeClient(c *Client) {
	if !c.in(s.Clients) {
		log.Print("Client is not part of the server")
		return
	}
	// leave all channels
	for _, ch := range s.Channels {
		ch.leave(c)
	}
	// leave server
	if i, ok := c.getIndex(s.Clients); ok {
		s.Clients = append(s.Clients[:i], s.Clients[i+1:]...)
	}
}

func (s *Server) getChannel(name string) *Channel {
	_, ok := s.Channels[name]
	if !ok {
		log.Print("Channel doesn't exist yet, creating a new one with the name ", name)
		s.Channels[name] = newChannel(name)
	}
	return s.Channels[name]
}
