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
}

// NewServer instantiates a server
func NewServer(cfg *Config) *Server {
	return &Server{
		Channels: make(map[string]*Channel),
		Clients:  make([]*Client, 0),
		Password: cfg.Password,
		Port:     cfg.Port,
	}
}

// Run starts the server
func (s *Server) Run() {
	ln, err := net.Listen("tcp", s.Port)
	defer ln.Close()
	check(err)
	log.Print("Listening on port ", s.Port)
	s.Listener = ln

	go s.startReceiving()
	// go s.startSending()
	s.acceptLoop()
}

func (s *Server) acceptLoop() {
	for {
		conn, err := s.Listener.Accept()
		check(err)
		log.Print("A new user connected with the remote IP: ", conn.RemoteAddr())
		client := newClient(conn)
		s.addClient(client)
	}
}

func (s *Server) startReceiving() {
	for {
		for _, v := range s.Clients {
			reader := bufio.NewReader(v.Conn)
			for {
				message, _ := reader.ReadString('\n')
				if len(message) == 0 {
					continue
				}
				message = strings.TrimSpace(message)
				//v.MsgQ <- message
				log.Print("***MSG: ", string(message))
				cmd, params := parseMessage(message)
				_, ok := commands[cmd]
				if ok {
					commands[cmd](params, v, s)
				}
			}
		}
	}
}

func (s *Server) startSending() {
	for {
		// TODO
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
	s.Clients = append(s.Clients, c)
}
