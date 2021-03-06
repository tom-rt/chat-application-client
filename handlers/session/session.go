package session

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"chat-application/client/handlers/connection"
	"chat-application/client/handlers/message"

	"github.com/gorilla/websocket"
)

func catchCtrlC(messageDialer *websocket.Conn, nickname string) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	for sig := range c {
		if sig == syscall.SIGINT {
			connection.Disconnect(messageDialer, nickname)
		}
	}
}

func readFromdialer(messageDialer *websocket.Conn) {
	for {
		_, message, err := messageDialer.ReadMessage()
		if err != nil {
			log.Println("Read error:", err)
			return
		}
		fmt.Printf("%s\n", message)
	}
}

func HandleSession(serverAddress string, nickname string) {
	messageUrl := url.URL{Scheme: "ws", Host: serverAddress, Path: "/run/session"}

	messageDialer, _, err := websocket.DefaultDialer.Dial(messageUrl.String(), nil)
	if err != nil {
		log.Fatal("Error dialing:", err)
	}
	defer messageDialer.Close()

	err = connection.Connect(messageDialer, nickname)

	if err != nil {
		log.Fatal("Cannot connect to server: ", err)
	}

	go catchCtrlC(messageDialer, nickname)

	go readFromdialer(messageDialer)

	stdio := make(chan string)
	go func() {
		for {
			var reader = bufio.NewReader(os.Stdin)
			message, _ := reader.ReadString('\n')

			stdio <- message
		}
	}()

	for msg := range stdio {
		payload := &message.MessageStruct{Connection: false, Disconnection: false, Nickname: nickname, Message: strings.Trim(msg, "\n")}
		marshaledPayload, err := json.Marshal(payload)
		if err != nil {
			fmt.Println("Error marshalling", err)
		}
		err = messageDialer.WriteMessage(websocket.TextMessage, []byte(string(marshaledPayload)))
		if err != nil {
			log.Println("Error writing message.", err)
			return
		}
	}
}
