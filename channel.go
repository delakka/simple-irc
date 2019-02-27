package main

// Channel represents a channel on an IRC server
type Channel struct {
	Name  string
	Users map[string]*User
	MsgQ  chan string
}

func newChannel(name string) *Channel {
	ch := &Channel{
		Name:  name,
		Users: make(map[string]*User),
		MsgQ:  make(chan string),
	}
	return ch
}

func (ch *Channel) sendLoop() {
	for message := range ch.MsgQ {
		for _, v := range ch.Users {
			v.MsgQ <- message
		}
	}
}

func (ch *Channel) send(message string) {
	ch.MsgQ <- message
}

func (ch *Channel) join(user *User) {
	ch.Users[user.Nick] = user
}

func (ch *Channel) leave(user *User) {
	_, ok := ch.Users[user.Nick]
	if ok {
		delete(ch.Users, user.Nick)
	}
}
