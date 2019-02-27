package main

// Channel represents a channel on an IRC server
type Channel struct {
	Name    string
	Clients map[string]*Client
	MsgQ    chan string
}

func newChannel(name string) *Channel {
	ch := &Channel{
		Name:    name,
		Clients: make(map[string]*Client),
		MsgQ:    make(chan string),
	}
	return ch
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
}

func (ch *Channel) join(client *Client) {
	ch.Clients[client.Nick] = client
}

func (ch *Channel) leave(client *Client) {
	_, ok := ch.Clients[client.Nick]
	if ok {
		delete(ch.Clients, client.Nick)
	}
}
