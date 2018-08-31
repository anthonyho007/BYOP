# Bring Your Own Playground (BYOP)

BYOP is a [Go](http://golang.org/) implementation of a host it yourself private chat server where you can chat and challenge your friends for a quick game.

### Status



### Installation

    go get github.com/gorilla/websocket

    go get github.com/anthonyho007/BYOP

### Usage

After cloning / downloading the repository, you can start the server by:

    go run main.go

then open http://localhost:8000 to access the chat server

### Chating

In the chat, you can use the following commands to perform different actions:

    Show the name of the current chatroom and all the users in the chatroom:
        /show
    
    List all available chatrooms:
        /list-chats

    List all online users:
        /online

    Enter a chatroom:
        /enter-chat <Chatroom_Name>

    Create a chatroom:
        /create-chat <Chatroom_Name>
    
    Invite friends into a private chatroom:
        /talk <username1> <username2> <username3>

