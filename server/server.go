package server

import (
	"fmt"

	"github.com/gorilla/websocket"
)

const (
	MaxNewConnections = 64
)

type Server struct {
	Id      string
	Clients map[string]*Client
	Chats   map[string]*Chat
	Hall    *Chat

	NewConnections chan *websocket.Conn
	ChangeChat     chan *Client
	Exit           chan int
}

func CreateServer() *Server {
	fmt.Println("Starting server")
	server := new(Server)
	go server.start()
	return server
}

func (server *Server) start() {
	server.Chats = make(map[string]*Chat)
	server.Clients = make(map[string]*Client)
	server.NewConnections = make(chan *websocket.Conn, MaxNewConnections)
	server.ChangeChat = make(chan *Client, MaxChatClientBuffer)
	server.Exit = make(chan int, 1)

	server.Hall = server.createChat(HallID, HallID)
	go server.Hall.start()
	go server.listen()

	<-server.Exit

}

func (server *Server) HandleNewConnection(conn *websocket.Conn) {
	server.NewConnections <- conn
}

func (server *Server) listen() {
	for {
		select {
		case conn := <-server.NewConnections:
			auth := Auth{}
			err := conn.ReadJSON(&auth)
			if err != nil {
				fmt.Println("connection err", err)
				continue
			}
			if auth.Name != "" && auth.Email != "" {
				client := server.getClient(auth.Name, auth.Email)
				if client == nil {
					fmt.Println("no client found create new one", auth.Name, auth.Email)
					client = server.createClient(conn, auth.Name, auth.Email)
				} else {
					client.Conn.Close()
				}
				client.Conn = conn
				go client.start()
			}

		case client := <-server.ChangeChat:
			if client.CurrentChat != nil {
				client.CurrentChat.LeaveChat <- client
			} else if client.ChangeChatId == "" {
				client.unblockRecvChannel()
			} else {
				for k := range client.Chats {
					fmt.Printf(k + " ")
				}
				chat := server.Chats[client.ChangeChatId]
				if chat == nil {
					fmt.Println("Failed to find chat id: " + client.ChangeChatId)
					continue
				}
				chat.EnterChat <- client
			}
		}
	}
}

func (server *Server) createClient(conn *websocket.Conn, name string, email string) *Client {
	client := createClient(conn, name, email)
	client.Server = server
	server.Clients[client.Id] = client
	client.ChangeChatId = HallID
	client.startChangingChat(HallID)
	return client
}

func (server *Server) createChat(id string, name string) *Chat {
	chat := createChat(id, name)
	server.Chats[chat.Id] = chat
	chat.Server = server
	return chat
}
