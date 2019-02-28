package main

import (
	"fmt"
	"log"
	"strings"
)

func check(err error) {
	if err != nil {
		log.Fatal("An error occured: ", err)
	}
}

func parseMessage(message string) (string, []string) {
	if len(message) != 0 {
		splitstr := strings.Split(message, " ")
		cmd := splitstr[0]
		if len(splitstr) > 1 {
			return cmd, splitstr[1:]
		}
		return cmd, []string{}
	}
	return "", []string{}
}

func buildMessage(prefix string, cmd string, params ...string) string {
	message := fmt.Sprintf(":%s %s", prefix, cmd)
	message += strings.Join(params, " ")
	return message
}
