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
		err := buildMessage(ERR_NONICKNAMEGIVEN.Numeric, c.Nick, ERR_NONICKNAMEGIVEN.Message)
		c.send(err)
		return
	}

	nick := params[0]
	if !s.isNickAvailable(nick) {
		err := buildMessage(s.Host, ERR_NICKNAMEINUSE.Numeric, "*", c.Nick, ERR_NICKNAMEINUSE.Message)
		c.send(err)
		return
	}
	c.Nick = nick
	log.Print("NICK set for client (", c.Conn.RemoteAddr(), "): ", nick)
}

func handleUSER(params []string, c *Client, s *Server) {
	if len(params) < 4 {
		err := buildMessage(ERR_NEEDMOREPARAMS.Numeric, c.Nick, ERR_NEEDMOREPARAMS.Message)
		c.send(err)
		return
	}

	c.UserName = params[0]
	c.HostName = params[2]
	c.RealName = strings.Join(params[3:], " ")

	// send RPL_WELCOME
	message := buildMessage(s.Host, RPL_WELCOME.Numeric, c.Nick, RPL_WELCOME.Message)
	c.send(message)
}

func handlePASS(params []string, c *Client, s *Server) {
	if s.auth(params[0]) {
		c.Authenticated = true
	}
	log.Print("Client ", c.Nick, " authenticated")
}

func handleJOIN(params []string, c *Client, s *Server) {
	if len(params) < 1 {
		err := buildMessage(ERR_NEEDMOREPARAMS.Numeric, c.Nick, ERR_NEEDMOREPARAMS.Message)
		c.send(err)
		return
	}

	channelName := params[0]
	if !strings.HasPrefix(channelName, "#") {
		channelName = "#" + channelName
	}
	ch := s.getChannel(channelName)
	ch.join(c)

	// send JOIN
	joinPrefix := fmt.Sprintf("%s!%s@%s", c.Nick, c.Nick, c.HostName)
	joinMessage := buildMessage(joinPrefix, "JOIN", ch.Name)
	for _, u := range ch.Clients {
		u.send(joinMessage)
	}

	// send RPL_TOPIC
	topicMessage := buildMessage(":"+s.Host, RPL_TOPIC.Numeric, c.Nick, ch.Name, ":"+ch.Topic)
	ch.MsgQ <- topicMessage
}

func handleNAMES(params []string, c *Client, s *Server) {
	for _, p := range params {
		for _, ch := range s.Channels {
			if p == ch.Name {
				clients := ch.getClientNicks()
				clientStr := strings.Join(clients, " ")
				log.Print("[D] clientlist: ", clientStr)
				message := buildMessage(":"+s.Host, RPL_NAMREPLY.Numeric, c.Nick, "=", ch.Name, clientStr)
				c.send(message)
			}
		}
	}
	endMessage := buildMessage(":"+s.Host, RPL_ENDOFNAMES.Numeric, c.Nick, RPL_ENDOFNAMES.Message)
	c.send(endMessage)
}

func handlePART(params []string, c *Client, s *Server) {
	if len(params) == 0 {
		err := buildMessage(ERR_NEEDMOREPARAMS.Numeric, c.Nick, ERR_NEEDMOREPARAMS.Message)
		c.send(err)
		return
	}

	chList := strings.Split(params[0], ",")
	suffix := ""
	if len(params) > 1 {
		suffix = strings.Join(params[1:], " ")
	}

	// send PART message
	for _, p := range chList {
		for _, ch := range c.Channels {
			if ch.Name == p {
				prefix := fmt.Sprintf(":%s!%s@%s", c.Nick, c.Nick, c.HostName)
				message := buildMessage(prefix, "PART", suffix)
				ch.MsgQ <- message
				ch.leave(c)
			}
		}
	}
}

func handlePRIVMSG(params []string, c *Client, s *Server) {
	if len(params) < 2 {
		err := buildMessage(ERR_NEEDMOREPARAMS.Numeric, c.Nick, ERR_NEEDMOREPARAMS.Message)
		c.send(err)
		return
	}

	// send PRIVSMG message(s)
	targetName := params[0]
	message := strings.Join(params[1:], " ")
	log.Print("Message to relay: ", message)

	if strings.HasPrefix(targetName, "#") {
		prefix := fmt.Sprintf(":%s!%s@%s", c.Nick, c.Nick, c.HostName)
		privMessage := buildMessage(prefix, "PRIVMSG", c.Channels[targetName].Name, message)
		c.Channels[targetName].MsgQ <- privMessage
	} else {
		if tar, ok := findClient(targetName, s.Clients); ok {
			prefix := fmt.Sprintf(":%s!%s@%s", c.Nick, c.Nick, c.HostName)
			privMessage := buildMessage(prefix, "PRIVMSG", tar.Nick, message)
			tar.send(privMessage)
		}
	}

}
