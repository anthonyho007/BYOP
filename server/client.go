package server

import (
	"fmt"
	"strings"
	"time"

	"github.com/gorilla/websocket"
)

const (
	MaxMessageBuffer = 36
)

type Client struct {
	Id   string
	Name string
	Code string

	Conn        *websocket.Conn
	Server      *Server
	CurrentChat *Chat

	Chats        map[string]*Chat
	ChangeChatId string
	ChangeChat   chan int

	SendMessage  chan Message
	CloseChannel chan int
}

func (client *Client) start() {
	fmt.Println("start client")
	go client.recv()
	go client.send()
	<-client.CloseChannel
}

func createClient(conn *websocket.Conn, name string, code string) *Client {
	client := &Client{

		Name: name,
		Code: code,

		Conn:        conn,
		CurrentChat: nil,
		Server:      nil,

		Chats: make(map[string]*Chat),

		SendMessage:  make(chan Message, MaxMessageBuffer),
		CloseChannel: make(chan int),
	}
	client.Id = generateId()
	return client
}

func (client *Client) generateMsg(str string) {
	msg := Message{
		Id:   "",
		Code: "",
		Name: "",
		Msg:  str,
		Date: time.Now(),
	}

	client.SendMessage <- msg
}

func (client *Client) recv() {
	for {
		select {
		case <-client.CloseChannel:
			goto EXIT
		default:
			{
			}
		}
		<-client.ChangeChat
		msg := Message{}
		err := client.Conn.ReadJSON(&msg)
		if err != nil {
			goto EXIT
		}
		msg.Date = time.Now()

		input := msg.Msg
		if strings.HasPrefix(input, "/") {
			if strings.HasPrefix(input, "/list-chats") {
				client.generateMsg("Chatrooms : " + client.Server.listAllChatnames())
				continue
			} else if strings.HasPrefix(input, "/show") {
				client.generateMsg(client.CurrentChat.getChatRoomInfo())
				continue

			} else if strings.HasPrefix(input, "/create-chat") {
				client.createChat(input)
				continue

			} else if strings.HasPrefix(input, "/enter-chat") {
				client.enterChat(input)
				continue

			} else if strings.HasPrefix(input, "/talk") {
				client.inviteToChat(input)
				continue
			} else if strings.HasPrefix(input, "/online") {
				client.generateMsg("Online : " + client.Server.listAllClientNames())
				continue
			}
		} else if input == "" {
			continue
		}

		client.CurrentChat.BroadcastMessage <- msg
	}
EXIT:
	client.Conn.Close()
}

func (client *Client) send() {
	for {
		select {
		case <-client.CloseChannel:
			goto EXIT
		case msg := <-client.SendMessage:
			err := client.Conn.WriteJSON(&msg)
			if err != nil {
				goto EXIT
			}
		}
	}

EXIT:
	client.Conn.Close()
}

func (client *Client) startChangingChat(chatId string) {
	client.ChangeChat = make(chan int)
	client.ChangeChatId = chatId
	client.Server.ChangeChat <- client
}

func (client *Client) unblockRecvChannel() {
	chatName := client.CurrentChat.Name
	client.generateMsg("You have entered chatroom " + chatName)
	client.ChangeChatId = ""
	close(client.ChangeChat)
}

func (client *Client) createChat(input string) string {
	name := strings.TrimPrefix(input, "/create-chat")
	name = strings.TrimSpace(name)
	chat := client.Server.createChat("", name)
	go chat.start()
	client.generateMsg("You have created chatroom " + name)
	return chat.Id
}

func (client *Client) enterChat(input string) {
	name := strings.TrimPrefix(input, "/enter-chat")
	name = strings.TrimSpace(name)
	chatId := client.Server.getChatId(name)
	if chatId == "" {
		chatId = client.createChat(name)
	}
	client.startChangingChat(chatId)
}

func (client *Client) inviteToChat(input string) {
	names := strings.TrimPrefix(input, "/talk")
	names = strings.TrimSpace(names)
	users := strings.Split(names, " ")

	clients := client.Server.getClientsByName(users)
	chatName := ""
	users = append(users, client.Name)
	for _, chat := range client.Server.Chats {
		if chat.containsOnlyClients(users) {
			chatName = chat.Id
		}
	}
	if chatName == HallID || chatName == "" {
		chatName = "room-" + generate3DigitId()
	}

	client.enterChat(chatName)

	for _, user := range clients {
		user.generateMsg(client.Name + " has invite you to chat in chatroom " + chatName)
	}

}
