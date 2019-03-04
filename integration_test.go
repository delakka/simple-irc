package main

import (
	"net"
	"strings"
	"testing"
)

func TestInit(t *testing.T) {
	cfg := NewConfig("./config.json")
	server := NewServer(cfg)
	go server.Run()
}

func sendTest(conn net.Conn, messages ...string) {
	for _, m := range messages {
		conn.Write([]byte(m + "\r\n"))
	}
}

func readTest(conn net.Conn) string {
	buffer := make([]byte, 1024)
	n, _ := conn.Read(buffer)

	messagesRaw := string(buffer[:n])
	messages := strings.Split(messagesRaw, "\r\n")
	return strings.TrimSpace(messages[0])
}

func prepareConnection() (net.Conn, error) {
	conn, err := net.Dial("tcp", ":6667")
	if err != nil {
		return nil, err
	}
	sendTest(conn, "PASS pw1!", "NICK user", "USER tstuser tstuser localhost :realname")

	in := make([]byte, 1024)
	conn.Read(in)
	conn.Read(in)

	return conn, nil
}

func TestPass_NoParams(t *testing.T) {
	conn, err := prepareConnection()
	if err != nil {
		t.Error("Failed to set up a working connection: ", err)
	}
	defer conn.Close()

	sendTest(conn, "PASS")

	in := readTest(conn)
	if in != "461 user :Need more parameters" {
		t.Error("Expected a different response. Received: ", in)
	}
}

func TestNick_NoParams(t *testing.T) {
	conn, err := prepareConnection()
	if err != nil {
		t.Error("Failed to set up a working connection: ", err)
	}
	defer conn.Close()

	sendTest(conn, "NICK")

	in := readTest(conn)
	if in != "431 user :No nick was given" {
		t.Error("Expected a different response. Received: ", in)
	}
}

func TestNick_InUse(t *testing.T) {
	conn, err := prepareConnection()
	if err != nil {
		t.Error("Failed to set up a working connection: ", err)
	}
	defer conn.Close()

	sendTest(conn, "NICK user")

	in := readTest(conn)
	if in != "irc.akka.io 433 * user :Nick is already in use" {
		t.Error("Expected a different response. Received: ", in)
	}
}

func TestUser_OneParam(t *testing.T) {
	conn, err := prepareConnection()
	if err != nil {
		t.Error("Failed to set up a working connection: ", err)
	}
	defer conn.Close()

	sendTest(conn, "USER hmm")

	in := readTest(conn)
	if in != "461 user :Need more parameters" {
		t.Error("Expected a different response. Received: ", in)
	}
}

func TestJoin_OneChannel(t *testing.T) {
	conn, err := prepareConnection()
	if err != nil {
		t.Error("Failed to set up a working connection: ", err)
	}
	defer conn.Close()

	sendTest(conn, "JOIN ch1")

	expected := []string{
		":user!user@localhost JOIN #ch1",
		":irc.akka.io 332 user #ch1 :TEST",
		":irc.akka.io 353 user = #ch1 user",
		":irc.akka.io 366 user #ch1 :End of NAMES list",
	}

	for _, e := range expected {
		in := readTest(conn)
		if in != e {
			t.Error("Expected a different response")
			t.FailNow()
		}
	}
}

func TestJoin_TwoChannels(t *testing.T) {
	conn, err := prepareConnection()
	if err != nil {
		t.Error("Failed to set up a working connection: ", err)
	}
	defer conn.Close()

	sendTest(conn, "JOIN ch11,ch22")

	expected := []string{
		":user!user@localhost JOIN #ch11",
		":irc.akka.io 332 user #ch11 :TEST",
		":irc.akka.io 353 user = #ch11 user",
		":irc.akka.io 366 user #ch11 :End of NAMES list",
		":user!user@localhost JOIN #ch22",
		":irc.akka.io 332 user #ch22 :TEST",
		":irc.akka.io 353 user = #ch22 user",
		":irc.akka.io 366 user #ch22 :End of NAMES list",
	}

	for _, e := range expected {
		in := readTest(conn)
		if in != e {
			t.Error("Expected a different response")
			t.FailNow()
		}
	}
}

func TestPart_OneChan(t *testing.T) {
	conn, err := prepareConnection()
	if err != nil {
		t.Error("Failed to set up a working connection: ", err)
	}
	defer conn.Close()
}
