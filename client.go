package main

import (
	"time"

	"github.com/gorilla/websocket"
)

const (
	MaxMessageBuffer = 36
)

type Client struct {
	Id    string
	Name  string
	Email string

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
	go client.recv()
	go client.send()
	<-client.CloseChannel
}

func createClient(conn *websocket.Conn, name string, email string) *Client {
	client := &Client{

		Name:  name,
		Email: email,

		Conn:        conn,
		CurrentChat: nil,
		Server:      nil,

		SendMessage:  make(chan Message, MaxMessageBuffer),
		CloseChannel: make(chan int),
	}
	client.Id = generateId()
	return client
}

func (client *Client) generateMsg(str string) {
	msg := Message{
		Id:    client.Id,
		Email: client.Email,
		Name:  client.Name,
		Msg:   str,
		Date:  time.Now(),
	}

	client.SendMessage <- msg
}

func (client *Client) recv() {
	for {
		select {
		case <-client.CloseChannel:
			goto EXIT
		}
		<-client.ChangeChat

		msg := Message{}
		err := client.Conn.ReadJSON(&msg)
		if err != nil {
			goto EXIT
		}
		msg.Date = time.Now()
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
	// client.Server.
}

func (client *Client) unblockRecvChannel() {
	client.ChangeChatId = ""
	close(client.ChangeChat)
}
