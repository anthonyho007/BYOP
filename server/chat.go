package server

import (
	"fmt"
	"strings"

	"github.com/anthonyho007/BYOP/datastructure"
)

const (
	HallID              = "Hall"
	MaxChatClientBuffer = 12
	MaxPastMessage      = 25
)

type Chat struct {
	Id string

	HistMessage *datastructure.CSlice

	Name    string
	Server  *Server
	Clients map[string]*Client

	BroadcastMessage chan Message
	EnterChat        chan *Client
	LeaveChat        chan *Client
}

func createChat(id string, name string) *Chat {
	chat := &Chat{
		Id:               id,
		Name:             name,
		Clients:          make(map[string]*Client),
		BroadcastMessage: make(chan Message, MaxMessageBuffer),
		EnterChat:        make(chan *Client, MaxChatClientBuffer),
		LeaveChat:        make(chan *Client, MaxChatClientBuffer),
	}
	if id == "" {
		chat.Id = generateId()
	}

	chat.HistMessage = datastructure.CSliceObj(MaxPastMessage)
	return chat
}

func (chat *Chat) start() {

	for {
		select {
		case client := <-chat.EnterChat:
			chat.enterChat(client)
			fmt.Println("print past message")
			for msg := range chat.HistMessage.List() {
				msg1 := msg.Entry.(Message)
				chat.BroadcastMessage <- msg1
			}
			client.unblockRecvChannel()
		case client := <-chat.LeaveChat:
			chat.leaveChat(client)
			chat.Server.ChangeChat <- client
		case msg := <-chat.BroadcastMessage:
			for _, client := range chat.Clients {
				client.SendMessage <- msg
			}
		}
	}
}

func (chat *Chat) enterChat(client *Client) {
	if client.CurrentChat != nil {
		return
	}
	fmt.Println("set current chat to " + chat.Name)
	chat.Clients[client.Id] = client
	client.CurrentChat = chat
	client.Chats[chat.Id] = chat
}

func (chat *Chat) leaveChat(client *Client) {
	if client.CurrentChat == nil {
		return
	}
	delete(chat.Clients, client.Id)
	client.CurrentChat = nil
}

func (chat *Chat) containsOnlyClients(clients []string) bool {
	if len(chat.Clients) != len(clients) {
		return false
	}
	for _, client := range chat.Server.getClientsByName(clients) {
		_, exist := chat.Clients[client.Id]
		if !exist {
			return false
		}
	}
	return true
}

func (chat *Chat) getAllChatMemberName() []string {
	var result []string
	for _, client := range chat.Clients {
		result = append(result, client.Name)
	}
	return result
}

func (chat *Chat) getChatRoomInfo() string {
	chatName := chat.Name
	members := strings.Join(chat.getAllChatMemberName(), ", ")
	msg := "You are currently at chatroom " + chatName + ", and " + members + " are currently in the chatroom"
	return msg
}
