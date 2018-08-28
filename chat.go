package main

const (
	HallID        = "Hall"
	MaxChatClient = 12
)

type Chat struct {
	Id string

	Server  *Server
	Clients map[string]*Client

	BroadcastMessage chan Message
	EnterChat        chan *Client
	LeaveChat        chan *Client
}

func createChat() *Chat {
	chat := &Chat{
		Clients:          make(map[string]*Client),
		BroadcastMessage: make(chan Message, MaxMessageBuffer),
		EnterChat:        make(chan *Client, MaxChatClient),
		LeaveChat:        make(chan *Client, MaxChatClient),
	}
	chat.Id = generateId()
	return chat
}

func (chat *Chat) enterChat(client *Client) {
	if client.CurrentChat != nil {
		return
	}
	chat.Clients[client.Id] = client
	client.CurrentChat = chat
	client.Chats[chat.Id] = chat
}

func (chat *Chat) leaveChat(client *Client) {
	if client.CurrentChat == nil {
		return
	}
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

func (chat *Chat) start() {
	for {

	}
}
