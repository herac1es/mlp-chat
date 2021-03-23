package server

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	uuid "github.com/satori/go.uuid"

	"github.com/gorilla/websocket"

	"github.com/herac1es/mlp-chat/pkg/safego"
)

var defaultManager = newManager()

type clientMap map[*Client]struct{}

// Manager: websocket server
type Manager struct {
	clients    clientMap
	broadcast  chan []byte
	register   chan *Client
	unregister chan *Client
}

func newManager() *Manager {
	return &Manager{
		clients:    make(clientMap),
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}
}

func (manager *Manager) start() {
	for {
		select {
		case conn := <-manager.register:
			manager.clients[conn] = struct{}{}
			bytes, _ := json.Marshal(Message{Content: fmt.Sprintf("%s has connected.", conn.id)})
			manager.send(bytes, clientMap{conn: struct{}{}})
		case conn := <-manager.unregister:
			if _, ok := manager.clients[conn]; ok {
				close(conn.send)
				delete(manager.clients, conn)
				bytes, _ := json.Marshal(Message{Content: fmt.Sprintf("%s has disconnected.", conn.id)})
				manager.send(bytes, clientMap{conn: struct{}{}})
			}
		case msg := <-manager.broadcast:
			manager.send(msg, nil)
		}
	}
}

// broadcast send message exclude all clients in ignore
func (manager *Manager) send(msg []byte, ignore clientMap) {
	for c := range manager.clients {
		if _, ok := ignore[c]; ok {
			continue
		}
		select {
		case c.send <- msg:
		default:
			close(c.send)
			delete(manager.clients, c)
		}
	}
}

func Run() {
	ctx := context.Background()
	fmt.Println("mlp-chat server run...")
	safego.Go(ctx, func(ctx context.Context) {
		defaultManager.start()
	})
	http.HandleFunc("/ws_chat", wsChat)
	http.ListenAndServe(":5268", nil)
}

func wsChat(writer http.ResponseWriter, request *http.Request) {
	conn, err := (&websocket.Upgrader{CheckOrigin: func(r *http.Request) bool { return true }}).Upgrade(writer, request, nil)
	if err != nil {
		fmt.Println(err)
		http.NotFound(writer, request)
		return
	}
	client := &Client{id: uuid.NewV4().String(), socket: conn, send: make(chan []byte)}

	defaultManager.register <- client
	safego.Go(request.Context(), func(ctx context.Context) {
		client.read()
	})
	safego.Go(request.Context(), func(ctx context.Context) {
		client.write()
	})
}
