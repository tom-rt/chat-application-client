package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"log"
	"net/url"
	"os"
	"strings"

	"github.com/gorilla/websocket"
)

type messageStruct struct {
	Connection    bool
	Disconnection bool
	Nickname      string
	Message       string
}

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

func handleSession(serverAddress string, nickname string) {
	messageUrl := url.URL{Scheme: "ws", Host: serverAddress, Path: "/run/session"}
	interrupt := make(chan os.Signal, 1)

	messageDialer, _, err := websocket.DefaultDialer.Dial(messageUrl.String(), nil)
	if err != nil {
		log.Fatal("Error dialing:", err)
	}
	defer messageDialer.Close()

	stdio := make(chan string)

	// CONNECTING
	connMessage := &messageStruct{Connection: true, Disconnection: false, Nickname: nickname, Message: "connecting !"}
	connJson, error := json.Marshal(connMessage)
	if error != nil {
		fmt.Println("Error while marshalling", error)
	}
	err = messageDialer.WriteMessage(websocket.TextMessage, []byte(string(connJson)))
	if err != nil {
		log.Println("Error writing message.", err)
		return
	}

	_, connectionResponse, err := messageDialer.ReadMessage()
	if err != nil {
		log.Println("Error on reading:", err)
		return
	}

	log.Printf("Connection reponse: %s", connectionResponse)

	// READING INCOMING MESSAGES FROM MESSAGE DIALER
	go func() {
		for {
			_, message, err := messageDialer.ReadMessage()
			if err != nil {
				log.Println("Read error:", err)
				return
			}
			log.Printf("Received: %s", message)
		}
	}()

	// READING STANDARD INPUT
	go func() {
		for {
			var reader = bufio.NewReader(os.Stdin)
			message, _ := reader.ReadString('\n')

			stdio <- message
		}
	}()

	for {
		select {
		case msg := <-stdio:

			payload := &messageStruct{Connection: false, Disconnection: false, Nickname: nickname, Message: strings.Trim(msg, "\n")}
			marshaledPayload, error := json.Marshal(payload)
			if error != nil {
				fmt.Println("error marshalling", error)
			}
			err = messageDialer.WriteMessage(websocket.TextMessage, []byte(string(marshaledPayload)))
			if err != nil {
				log.Println("Error writing message.", err)
				return
			}

		case <-interrupt:
			err = messageDialer.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, "Closing client"))
			if err != nil {
				log.Println("Error on closing client:", err)
			}
			return
		}
	}
}

func main() {
	nickname, serverAddress, err := parseFlags()

	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println("your nickname is:", nickname)
	fmt.Println("server address:", serverAddress)

	handleSession(serverAddress, nickname)
}
