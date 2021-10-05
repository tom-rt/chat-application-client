package connection

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/gorilla/websocket"

	"chat-application/client/handlers/message"
)

type ConnectionResponse struct {
	IsAllowed bool
	Message   string
}

func Connect(messageDialer *websocket.Conn, nickname string) {
	// CONNECTING
	connMessage := &message.MessageStruct{Connection: true, Disconnection: false, Nickname: nickname, Message: "connecting !"}
	connJson, error := json.Marshal(connMessage)
	if error != nil {
		log.Println("Error while marshalling:", error)
	}
	err := messageDialer.WriteMessage(websocket.TextMessage, []byte(string(connJson)))
	if err != nil {
		log.Println("Error writing message:", err)
		return
	}

	_, resp, err := messageDialer.ReadMessage()
	if err != nil {
		log.Println("Error on reading:", err)
		return
	}

	var connectionResponse ConnectionResponse
	json.Unmarshal(resp, &connectionResponse)

	fmt.Printf("%s\n", connectionResponse.Message)

	if !connectionResponse.IsAllowed {
		err = messageDialer.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, "Closing client"))
		if err != nil {
			log.Println("Error on closing client:", err)
			return
		}
		os.Exit(0)
	}
}

func Disconnect(messageDialer *websocket.Conn, nickname string) {
	connMessage := &message.MessageStruct{Connection: false, Disconnection: true, Nickname: nickname, Message: "disconnecting"}
	connJson, error := json.Marshal(connMessage)
	if error != nil {
		log.Println("Error while marshalling:", error)
	}
	err := messageDialer.WriteMessage(websocket.TextMessage, []byte(string(connJson)))
	if err != nil {
		log.Println("Error on closing client:", err)
	}

	err = messageDialer.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, "Closing client"))
	if err != nil {
		log.Println("Error on closing client:", err)
	}
	os.Exit(0)
}
