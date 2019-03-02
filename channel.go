package main

import "log"

// Channel represents a channel on an IRC server
type Channel struct {
	Name    string
	Topic   string
	Clients map[string]*Client
	MsgQ    chan string
}

func newChannel(name string) *Channel {
	ch := &Channel{
		Name:    name,
		Topic:   "TEST",
		Clients: make(map[string]*Client),
		MsgQ:    make(chan string),
	}
	return ch
}

func (ch *Channel) setTopic(topic string) {
	ch.Topic = topic
}

func (ch *Channel) sendLoop() {
	for message := range ch.MsgQ {
		for _, v := range ch.Clients {
			v.send(message)
		}
	}
}

func (ch *Channel) send(message string) {
	ch.MsgQ <- message
	// for _, v := range ch.Clients {
	// 	v.send(message)
	// }
}

func (ch *Channel) join(client *Client) {
	ch.Clients[client.Nick] = client
	client.joinChannel(ch)

	log.Print("Added client ", client.Nick, " to channel ", ch.Name)
}

func (ch *Channel) leave(client *Client) {
	if _, ok := ch.Clients[client.Nick]; ok {
		delete(ch.Clients, client.Nick)
	}
	client.leaveChannel(ch)
}

func (ch *Channel) getClientNicks() []string {
	clients := make([]string, 0)
	for _, c := range ch.Clients {
		clients = append(clients, c.Nick)
	}
	return clients
}
