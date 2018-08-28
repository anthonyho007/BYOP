package main

import (
	"fmt"
	"math/rand"

	"github.com/gorilla/websocket"
)

type Server struct {
	Id      string
	Clients map[string]*Client
	Chat    map[string]*Chat
	Hall    *Chat
}

func (server *Server) createClient(conn *websocket.Conn, name string, email string) *Client {
	client := createClient(conn, name, email)
	server.Clients[client.Id] = client
	return client
}

func (server *Server) getClient(name string, email string) *Client {
	for _, client := range server.Clients {
		if client.Name == name && client.Email == email {
			return client
		}
	}
	return nil
}

func (server *Server) getClientByName(name string) *Client {
	for _, client := range server.Clients {
		if client.Name == name {
			return client
		}
	}
	return nil
}

func (server *Server) getClientsByName(names []string) []*Client {
	clients := []*Client{}
	for _, name := range names {
		client := server.getClientByName(name)
		if client != nil {
			clients = append(clients, client)
		}
	}
	return clients
}

func generateId() string {
	id := fmt.Sprintf("%d", rand.Intn(999999))
	return id
}
