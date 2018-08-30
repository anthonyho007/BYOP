package server

const (
	HallID              = "Hall"
	MaxChatClientBuffer = 12
)

type Chat struct {
	Id string

	Server  *Server
	Clients map[string]*Client

	BroadcastMessage chan Message
	EnterChat        chan *Client
	LeaveChat        chan *Client
}

func createChat(id string) *Chat {
	chat := &Chat{
		Id:               id,
		Clients:          make(map[string]*Client),
		BroadcastMessage: make(chan Message, MaxMessageBuffer),
		EnterChat:        make(chan *Client, MaxChatClientBuffer),
		LeaveChat:        make(chan *Client, MaxChatClientBuffer),
	}
	if id == "" {
		chat.Id = generateId()
	}
	return chat
}

func (chat *Chat) start() {

	for {
		select {
		case client := <-chat.EnterChat:
			chat.enterChat(client)
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
