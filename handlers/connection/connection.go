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

func Connect(messageDialer *websocket.Conn, nickname string) error {
	// CONNECTING
	connMessage := &message.MessageStruct{Connection: true, Disconnection: false, Nickname: nickname, Message: "connecting !"}
	connJson, err := json.Marshal(connMessage)
	if err != nil {
		log.Println("Error while marshalling:", err)
		return err
	}
	err = messageDialer.WriteMessage(websocket.TextMessage, []byte(string(connJson)))
	if err != nil {
		log.Println("Error writing message:", err)
		return err
	}

	_, resp, err := messageDialer.ReadMessage()
	if err != nil {
		log.Println("Error on reading:", err)
		return err
	}

	var connectionResponse ConnectionResponse
	err = json.Unmarshal(resp, &connectionResponse)
	if err != nil {
		log.Println("Unmarshaling error: ", err)
		return err
	}

	fmt.Printf("%s\n", connectionResponse.Message)

	if !connectionResponse.IsAllowed {
		err = messageDialer.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, "Closing client"))
		if err != nil {
			log.Println("Error on closing client:", err)
			return err
		}
	}
	return nil
}

func Disconnect(messageDialer *websocket.Conn, nickname string) {
	connMessage := &message.MessageStruct{Connection: false, Disconnection: true, Nickname: nickname, Message: "disconnecting"}
	connJson, err := json.Marshal(connMessage)
	if err != nil {
		log.Println("Error while marshalling:", err)
	}
	err = messageDialer.WriteMessage(websocket.TextMessage, []byte(string(connJson)))
	if err != nil {
		log.Println("Error on closing client:", err)
	}

	err = messageDialer.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, "Closing client"))
	if err != nil {
		log.Println("Error on closing client:", err)
	}
	os.Exit(0)
}
