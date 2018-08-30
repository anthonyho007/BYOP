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
	fmt.Printf("create serverr")
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
	fmt.Println("create new hall and server")
	server.Hall = server.createChat(HallID)
	go server.Hall.start()
	go server.listen()

	<-server.Exit

}

func (server *Server) HandleNewConnection(conn *websocket.Conn) {
	fmt.Println("send conn to newconnection chan")
	// auth := Auth{}
	// fmt.Print("okay")
	// conn.ReadJSON(&auth)
	// fmt.Println(auth.Name)
	server.NewConnections <- conn
}

func (server *Server) listen() {
	for {
		fmt.Println("loop on server listen")
		select {
		case conn := <-server.NewConnections:
			auth := Auth{}
			fmt.Println("stuck here")
			err := conn.ReadJSON(&auth)
			if err != nil {
				fmt.Println("conn err", err)
			}
			fmt.Println("got auth response")
			fmt.Println(auth.Name)
			if auth.Name != "" && auth.Email != "" {
				client := server.getClient(auth.Name, auth.Email)
				if client == nil {
					fmt.Println("no client found create new one", auth.Name, auth.Email)
					// client = server.createClient(conn, auth.Name, auth.Email)
					client = server.createClient(conn, auth.Name, auth.Email)
					// server.Clients[client.Id] = client
					// client.ChangeChatId = HallID
					// client.startChangingChat(HallID)
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
				fmt.Println("client changechat id and current chat room", client.ChangeChatId)
				for k := range client.Chats {
					fmt.Printf(k + " ")
				}
				chat := server.Chats[client.ChangeChatId]
				if chat == nil {
					fmt.Println("create a new chat")
					chat = server.createChat("")
					go chat.start()
				}
				fmt.Println("entering chat " + chat.Id)
				chat.EnterChat <- client
			}
		}
	}
}

func (server *Server) createClient(conn *websocket.Conn, name string, email string) *Client {
	fmt.Println("create client")
	client := createClient(conn, name, email)
	client.Server = server
	server.Clients[client.Id] = client
	client.ChangeChatId = HallID
	fmt.Println("HallId ", HallID)
	client.startChangingChat(HallID)
	return client
}

func (server *Server) createChat(id string) *Chat {
	if id == "" {
		id = generateId()
	}
	fmt.Println("Print generated Id", id)
	chat := createChat(id)
	server.Chats[chat.Id] = chat
	return chat
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
