package main

import (
	"log"
	"strings"
)

func check(err error) {
	if err != nil {
		log.Fatal("An error occured: ", err)
	}
}

func parseMessage(message string) (string, []string) {
	splitstr := strings.Split(message, " ")
	return splitstr[0], splitstr[1:]
}
