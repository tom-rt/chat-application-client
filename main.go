package main

import (
	"errors"
	"flag"
	"fmt"
)

func parseFlags() (string, string, error) {
	var nickname string
	var host string
	var port string
	var helpMessage string = "A nickname argument is required. It must be less than 20 characters long, exiting."

	flag.StringVar(&nickname, "nickname", "", helpMessage)
	flag.StringVar(&host, "host", "localhost", helpMessage)
	flag.StringVar(&port, "port", "8080", helpMessage)
	flag.Parse()

	if nickname == "" || len(nickname) > 20 {
		fmt.Println(helpMessage)
		return "", "", errors.New("no nickname provided")
	}

	return nickname, host + ":" + port, nil
}

func main() {
	nickname, serverAddress, err := parseFlags()

	if err != nil {
		return
	}
	fmt.Println("your nickname is:", nickname)
	fmt.Println("server address:", serverAddress)

}
