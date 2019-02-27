package main

import "log"

type commandHandler func([]string, *Client, *Server)

var commands = map[string]commandHandler{
	"NICK":    handleNICK,
	"USER":    handleUSER,
	"PASS":    handlePASS,
	"JOIN":    handleJOIN,
	"NAMES":   handleNAMES,
	"PART":    handlePART,
	"PRIVMSG": handlePRIVMSG,
}

func handleNICK(params []string, cl *Client, s *Server) {
	if len(params) != 1 {
		log.Print("NICK's param count is not 1")
		cl.send("There should only be 1 parameter")
		return
	}

	nick := params[0]
	if !s.isNickAvailable(nick) {
		log.Print("Nick already in use")
		cl.send("Nick already in use")
		return
	}
	cl.Nick = nick
	cl.send("Welcome, " + nick)
}

func handleUSER(params []string, cl *Client, s *Server) {}

func handlePASS(params []string, cl *Client, s *Server) {
	if s.auth(params[0]) {
		cl.Authenticated = true
	}
	cl.send("Authenticated.")
}

func handleJOIN(params []string, cl *Client, s *Server) {}

func handleNAMES(params []string, cl *Client, s *Server) {}

func handlePART(params []string, cl *Client, s *Server) {}

func handlePRIVMSG(params []string, cl *Client, s *Server) {}
