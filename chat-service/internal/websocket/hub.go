package websocket

import (
	"encoding/json"
	"fmt"

	"chat-service/config"
	grpcClients "chat-service/internal/gRPC/clients"
)

// Hub maintains the set of active clients and broadcasts messages to the
// clients.
type Hub struct {
	// Registered clients.
	clients map[*Client]bool

	// Inbound messages from the clients.
	broadcast chan []byte

	// Register requests from the clients.
	register chan *Client

	// Unregister requests from clients.
	unregister chan *Client
}

type Message struct {
	Message      string `json:"message"`
	SenderID     string `json:"senderId"`
	RecevierType string `json:"recevierType"`
	RecevierId   string `json:"recevierId"`
}

//Receiver types Enum
const (
	GROUP = "group"
	PERSONEL = "personel"
)

func NewHub() *Hub {
	return &Hub{
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			registerNewConnection(h, client)
		case client := <-h.unregister:
			unregisterConnection(h, client)
		case message := <-h.broadcast:
			handleBroadcast(h, message)
		}
	}
}

func registerNewConnection(h *Hub, client *Client) {
	h.clients[client] = true
	logger.Infof("Client %s registered\n", client.conn.RemoteAddr().String())
	logger.Infof("Number of connected clients: %d\n", len(h.clients))

	websocketManagerClient, err := grpcClients.NewWebsocketManagerClient(config.Env.WebsocketManagerUrl)
	defer websocketManagerClient.Disconnect()
	if err != nil {
		logger.Errorf("Error creating websocket manager client: %v\n", err)
		return
	}

	data := map[string]string{
		"connection": client.conn.RemoteAddr().String(),
		"server":     fmt.Sprintf("%s:%s", config.Env.Host, config.Env.Port),
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		logger.Errorf("Error marshalling data: %v\n", err)
		return
	}

	err = websocketManagerClient.Register(client.userId, string(jsonData))
	if err != nil {
		logger.Errorf("Error registering client with websocket manager: %v\n", err)
	}
}

func unregisterConnection(h *Hub, client *Client) {
	if _, ok := h.clients[client]; !ok {
		logger.Errorf("Client %s not found\n", client.conn.RemoteAddr().String())
		return
	}

	delete(h.clients, client)
	close(client.send)

	websocketManagerClient, err := grpcClients.NewWebsocketManagerClient(config.Env.WebsocketManagerUrl)
	defer websocketManagerClient.Disconnect()
	if err != nil {
		logger.Errorf("Error creating websocket manager client: %v\n", err)
		return
	}

	err = websocketManagerClient.UnRegister(client.userId)
	if err != nil {
		logger.Errorf("Error unregistering client with websocket manager: %v\n", err)
	}
}

func handleBroadcast(h *Hub, message []byte) {
	// if the connection is on the current server publish to him
	// if not then send it to the message service or group message service
	logger.Infof("Broadcasting message: %s\n", string(message))
	data := Message{}
	err := json.Unmarshal(message, &data)
	if err != nil {
		logger.Errorf("Error unmarshalling message: %v\n", err)
		return
	}

	switch data.RecevierType {
	case GROUP:
		// send message to group group message service
	case PERSONEL:
		// send message to message service
	default:
		logger.Errorf("Invalid Message recevier type: %s\n", data.RecevierType)
		return
	}

	for client := range h.clients {
		select {
		case client.send <- message:
		default:
			close(client.send)
			delete(h.clients, client)
		}
	}
}
