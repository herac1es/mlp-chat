package server

import (
	"encoding/json"
	"fmt"

	"github.com/gorilla/websocket"
)

// Client: every instance present a websocket connection
type Client struct {
	id     string
	socket *websocket.Conn
	send   chan []byte
}

func (client *Client) read() {
	defer func() {
		defaultManager.unregister <- client
		client.socket.Close()
	}()

	for {
		_, message, err := client.socket.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				fmt.Println(err)
			}
			break
		}
		bytes, _ := json.Marshal(Message{
			Sender:  client.id,
			Content: string(message),
		})
		defaultManager.broadcast <- bytes
	}
}

func (client *Client) write() {
	defer func() {
		client.socket.Close()
	}()
	for {
		select {
		case msg, ok := <-client.send:
			if !ok {
				client.socket.WriteMessage(websocket.CloseMessage, []byte{})
			}
			client.socket.WriteMessage(websocket.TextMessage, msg)
		}
	}
}
