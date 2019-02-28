package main

import (
	"fmt"
	"log"
	"strings"
)

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

func handleNICK(params []string, c *Client, s *Server) {
	if len(params) == 0 {
		//ERR_NONICKNAMEGIVEN
		return
	}

	nick := params[0]
	if !s.isNickAvailable(nick) {
		//ERR_NICKNAMEINUSE
		return
	}
	c.Nick = nick
}

func handleUSER(params []string, c *Client, s *Server) {
	if len(params) < 4 {
		//ERR_NEEDMOREPARAMS
		return
	}

	c.UserName = params[0]
	c.HostName = params[2]
	c.RealName = params[3]

	// send RPL_WELCOME
	msg := fmt.Sprintf(":%s %s", s.Host, RPL_WELCOME)
	msg = msg + " " + c.UserName + ": Welcome to the server!"
	log.Print(msg)
	c.send(msg)
}

func handlePASS(params []string, c *Client, s *Server) {
	if s.auth(params[0]) {
		c.Authenticated = true
	}
	c.send("Authenticated.")
}

func handleJOIN(params []string, c *Client, s *Server) {
	chName := params[0]
	if !strings.HasPrefix(chName, "#") {
		chName = "#" + chName
	}
	ch := s.getChannel(chName)
	ch.join(c)

	msg_join := fmt.Sprintf("%s!~%s@%s", c.UserName, c.UserName, s.Host)
	msg_join = msg_join + " JOIN " + ch.Name
	c.send(msg_join)
}

func handleNAMES(params []string, c *Client, s *Server) {}

func handlePART(params []string, c *Client, s *Server) {}

func handlePRIVMSG(params []string, c *Client, s *Server) {}
