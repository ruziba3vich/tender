package websocket

import (
	"encoding/json"
	"sync"

	"github.com/gorilla/websocket"
)

type Client struct {
	ID       string
	Conn     *websocket.Conn
	TenderID string
	mu       sync.Mutex
}

type Manager struct {
	clients    map[*Client]bool
	register   chan *Client
	unregister chan *Client
	broadcast  chan []byte
	mu         sync.RWMutex
}

func NewManager() *Manager {
	return &Manager{
		clients:    make(map[*Client]bool),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		broadcast:  make(chan []byte),
	}
}

// RegisterClient registers a new client
func (m *Manager) RegisterClient(client *Client) {
	m.register <- client
}

// UnregisterClient unregisters a client
func (m *Manager) UnregisterClient(client *Client) {
	m.unregister <- client
}

func (m *Manager) Run() {
	for {
		select {
		case client := <-m.register:
			m.mu.Lock()
			m.clients[client] = true
			m.mu.Unlock()

		case client := <-m.unregister:
			m.mu.Lock()
			if _, ok := m.clients[client]; ok {
				delete(m.clients, client)
				client.Conn.Close()
			}
			m.mu.Unlock()

		case message := <-m.broadcast:
			m.broadcastMessage(message)
		}
	}
}

func (m *Manager) BroadcastToTender(tenderID string, message interface{}) error {
	data, err := json.Marshal(message)
	if err != nil {
		return err
	}

	m.mu.RLock()
	defer m.mu.RUnlock()

	for client := range m.clients {
		if client.TenderID == tenderID {
			client.mu.Lock()
			err := client.Conn.WriteMessage(websocket.TextMessage, data)
			if err != nil {
				client.Conn.Close()
				delete(m.clients, client)
			}
			client.mu.Unlock()
		}
	}
	return nil
}

func (m *Manager) broadcastMessage(message []byte) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	for client := range m.clients {
		client.mu.Lock()
		err := client.Conn.WriteMessage(websocket.TextMessage, message)
		if err != nil {
			client.Conn.Close()
			delete(m.clients, client)
		}
		client.mu.Unlock()
	}
}
