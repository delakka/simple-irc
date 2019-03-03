package main

import (
	"fmt"
	"log"
	"strings"
)

type commandHandler func([]string, *Client, *Server)

// A map of the commands currently implemented and their handlers
var commands = map[string]commandHandler{
	"NICK":    handleNICK,
	"USER":    handleUSER,
	"PASS":    handlePASS,
	"JOIN":    handleJOIN,
	"NAMES":   handleNAMES,
	"PART":    handlePART,
	"PRIVMSG": handlePRIVMSG,
	"QUIT":    handleQUIT,
}

func handlePASS(params []string, c *Client, s *Server) {
	if s.auth(params[0]) {
		c.Authenticated = true
	}
	log.Print("Client ", c.Nick, " authenticated")
}

func handleNICK(params []string, c *Client, s *Server) {
	if !c.Authenticated {
		return
	}
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

	// send NICK
	message := buildMessage(fmt.Sprintf(":%s!%s@%s", c.Nick, c.Nick, c.HostName), "NICK", nick)
	c.send(message)

	c.setNick(nick)

	log.Print("NICK set for client (", c.Conn.RemoteAddr(), "): ", nick)
}

func handleUSER(params []string, c *Client, s *Server) {
	if !c.Authenticated {
		return
	}
	if len(params) < 4 {
		err := buildMessage(ERR_NEEDMOREPARAMS.Numeric, c.Nick, ERR_NEEDMOREPARAMS.Message)
		c.send(err)
		return
	}

	c.UserName = params[0]
	c.HostName = params[2]
	c.RealName = strings.Join(params[3:], " ")
	c.Registered = true

	// send RPL_WELCOME
	message := buildMessage(s.Host, RPL_WELCOME.Numeric, c.Nick, RPL_WELCOME.Message)
	c.send(message)
}

func handleJOIN(params []string, c *Client, s *Server) {
	if !c.Registered {
		return
	}
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
	joinPrefix := fmt.Sprintf(":%s!%s@%s", c.Nick, c.Nick, c.HostName)
	joinMessage := buildMessage(joinPrefix, "JOIN", ch.Name)
	for _, u := range ch.Clients {
		u.send(joinMessage)
	}

	// send RPL_TOPIC
	topicMessage := buildMessage(fmt.Sprintf(":%s", s.Host), RPL_TOPIC.Numeric, c.Nick, ch.Name, fmt.Sprintf(":%s", ch.Topic))
	c.send(topicMessage)
	sendNamesToChannel(ch, c, s)
}

func handleNAMES(params []string, c *Client, s *Server) {
	if !c.Registered {
		return
	}
	for _, p := range params {
		for _, ch := range s.Channels {
			if p == ch.Name {
				sendNamesToChannel(ch, c, s)
			}
		}
	}
}

func sendNamesToChannel(ch *Channel, c *Client, s *Server) {
	clients := ch.getClientNicks()
	clientStr := strings.Join(clients, " ")

	message := buildMessage(fmt.Sprintf(":%s", s.Host), RPL_NAMREPLY.Numeric, c.Nick, "=", ch.Name, clientStr)
	c.send(message)

	endMessage := buildMessage(fmt.Sprintf(":%s", s.Host), RPL_ENDOFNAMES.Numeric, c.Nick, ch.Name, RPL_ENDOFNAMES.Message)
	c.send(endMessage)
}

func handlePART(params []string, c *Client, s *Server) {
	if !c.Registered {
		return
	}
	if len(params) == 0 {
		err := buildMessage(ERR_NEEDMOREPARAMS.Numeric, c.Nick, ERR_NEEDMOREPARAMS.Message)
		c.send(err)
		return
	}

	chList := strings.Split(params[0], ",")
	suffix := "bye!"
	if len(params) > 1 {
		suffix = strings.Join(params[1:], " ")
	}

	// send PART message
	for _, p := range chList {
		for _, ch := range c.Channels {
			if ch.Name == p {
				prefix := fmt.Sprintf(":%s!%s@%s", c.Nick, c.Nick, c.HostName)
				message := buildMessage(prefix, "PART", ch.Name, suffix)
				c.send(message)
				ch.leave(c)
			}
		}
	}
}

func handlePRIVMSG(params []string, c *Client, s *Server) {
	if !c.Registered {
		return
	}
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
		if ch, ok := c.Channels[targetName]; ok {
			prefix := fmt.Sprintf(":%s!%s@%s", c.Nick, c.Nick, c.HostName)
			privMessage := buildMessage(prefix, "PRIVMSG", ch.Name, message)
			ch.send(c, privMessage)
		} else {
			err := buildMessage(ERR_NORECIPIENT.Numeric, c.Nick, ERR_NORECIPIENT.Message)
			c.send(err)
		}
	} else {
		if tar, ok := findClient(targetName, s.Clients); ok {
			log.Print("target found: ", tar.Nick)
			prefix := fmt.Sprintf(":%s!%s@%s", c.Nick, c.Nick, c.HostName)
			privMessage := buildMessage(prefix, "PRIVMSG", tar.Nick, message)
			tar.send(privMessage)
		} else {
			err := buildMessage(ERR_NORECIPIENT.Numeric, c.Nick, ERR_NORECIPIENT.Message)
			c.send(err)
		}
	}

}

func handleQUIT(params []string, c *Client, s *Server) {
	log.Print("Closing the connection for client ", c.Nick, c.Conn.RemoteAddr().String())
	s.removeClient(c)
	c.Conn.Close()
}
