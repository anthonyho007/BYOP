package server

import (
	"fmt"
	"math/rand"
	"strings"
)

func generateId() string {
	id := fmt.Sprintf("%d", 100+rand.Intn(999999))
	return id
}

func generate3DigitId() string {
	id := fmt.Sprintf("%d", 100+rand.Intn(899))
	return id
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

func (server *Server) getAllChatName() []string {
	chatNames := []string{}

	for _, chat := range server.Chats {
		chatNames = append(chatNames, chat.Name)
	}
	return chatNames
}

func (server *Server) getAllClientName() []string {
	clientNames := []string{}

	for _, client := range server.Clients {
		clientNames = append(clientNames, client.Name)
	}
	return clientNames
}

func (server *Server) listAllChatnames() string {
	return strings.Join(server.getAllChatName(), ", ")
}

func (server *Server) listAllClientNames() string {
	return strings.Join(server.getAllClientName(), ", ")
}

func (server *Server) getChatId(name string) string {
	id := ""
	for _, chat := range server.Chats {
		if chat.Name == name {
			return chat.Id
		}
	}
	return id
}
