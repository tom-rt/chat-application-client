package main

import (
	"errors"
	"flag"
	"fmt"
	"log"

	"chat-application/client/handlers/session"
)

func parseFlags() (string, string, error) {
	var nickname string
	var host string
	var port string
	var helpMessage string = "A nickname argument is required. It must be less than 20 characters long, exiting."

	flag.StringVar(&nickname, "nickname", "", helpMessage)
	flag.StringVar(&host, "host", "localhost", "invalid host value")
	flag.StringVar(&port, "port", "8080", "invalid port value")
	flag.Parse()

	if nickname == "" || len(nickname) > 20 {
		fmt.Println(helpMessage)
		return "", "", errors.New("invalid nickname")
	}

	return nickname, host + ":" + port, nil
}

func main() {
	nickname, serverAddress, err := parseFlags()

	if err != nil {
		log.Println("Error:", err)
		return
	}

	session.HandleSession(serverAddress, nickname)
}
