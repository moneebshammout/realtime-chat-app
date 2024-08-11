package websocket

import (
	"encoding/json"
	"fmt"

	appConfig "chat-service/config/app"
	"chat-service/config/queues"

	grpcClients "chat-service/internal/gRPC/clients"
)

// Hub maintains the set of active clients and broadcasts messages to the
// clients.
type Hub struct {
	// Registered clients.
	clients map[*Client]bool

	// Inbound messages from the clients to send to message or group message services.
	broadcast chan []byte

	// Register requests from the clients.
	register chan *Client

	// Unregister requests from clients.
	unregister chan *Client

	// send message to client
	Send chan *SendMessage
}

type Message struct {
	Message      string `json:"message"`
	SenderID     string `json:"senderId"`
	RecevierType string `json:"recevierType"`
	RecevierId   string `json:"recevierId"`
}

type SendMessage struct {
	Message    string `json:"message"`
	SenderId   string `json:"senderId"`
	RecevierId string `json:"recevierId"`
}

// Receiver types Enum
const (
	GROUP    = "group"
	PERSONEL = "personel"
)

func NewHub() *Hub {
	return &Hub{
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
		Send:  make(chan *SendMessage),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			go registerNewConnection(h, client)
		case client := <-h.unregister:
			go unregisterConnection(h, client)
		case message := <-h.broadcast:
			go handleBroadcast(h, message)

		case message := <-h.Send:
			go handlePublishMessage(h, message)
		}
	}
}

func registerNewConnection(h *Hub, client *Client) {
	h.clients[client] = true
	logger.Infof("Client %s registered\n", client.conn.RemoteAddr().String())
	logger.Infof("Number of connected clients: %d\n", len(h.clients))

	websocketManagerClient, err := grpcClients.NewWebsocketManagerClient(appConfig.Env.WebsocketManagerUrl)
	defer websocketManagerClient.Disconnect()
	if err != nil {
		logger.Errorf("Error creating websocket manager client: %v\n", err)
		return
	}

	data := map[string]string{
		"connection": client.conn.RemoteAddr().String(),
		"server":     fmt.Sprintf("%s:%s", appConfig.Env.Host, appConfig.Env.Port),
		"messageQueue": queues.Env.MessageQueue,
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

	websocketManagerClient, err := grpcClients.NewWebsocketManagerClient(appConfig.Env.WebsocketManagerUrl)
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
		handleGroupMessage(h, data)
	case PERSONEL:
		handlePersonelMessage(h, data)
	default:
		logger.Errorf("Invalid Message recevier type: %s\n", data.RecevierType)
	}
}

func handlePersonelMessage(h *Hub, message Message) {
	// send message to message service
	var found bool
	for client := range h.clients {
		if client.userId == message.RecevierId {
			select {
			case client.send <- []byte(message.Message):
			default:
				close(client.send)
				delete(h.clients, client)
			}

			found = true
			break
		}
	}

	if found {
		// return
	}

	logger.Infof("Client not found sending to message service: %s\n", message.RecevierId)
	messageClient, err := grpcClients.NewMessageServiceClient(appConfig.Env.MessageServiceUrl)
	defer messageClient.Disconnect()
	if err != nil {
		logger.Errorf("Error creating message client: %v\n", err)
		return
	}

	data := map[string]string{
		"message":    message.Message,
		"senderId":   message.SenderID,
		"receiverId": message.RecevierId,
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		logger.Errorf("Error marshalling message: %v\n", err)
		return
	}

	err = messageClient.Send(string(jsonData))
	if err != nil {
		logger.Errorf("Error sending message to message service: %v\n", err)
		return
	}

	logger.Infof("Message sent to message service: %s\n", message.Message)
}

func handleGroupMessage(h *Hub, message Message) {
	// send message to group group message service
}


func handlePublishMessage(h *Hub, message *SendMessage) {
	var found bool
	for client := range h.clients {
		if client.userId == message.RecevierId {
			select {
			case client.send <- []byte(message.Message):
			default:
				close(client.send)
				delete(h.clients, client)
			}

			found = true
			break
		}
	}

	if found {
		logger.Info("Message Published successfully")
		return
	}else{
		//disconnected from server during sending message send to relay service
		logger.Error("Client not found")
	}


}