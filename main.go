package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/anthonyho007/BYOP/server"
	"github.com/gorilla/websocket"
)

func main() {
	// acquire port
	port := os.Getenv("PORT")

	if port == "" {
		log.Fatal("$PORT value must be define")
	}
	wsServer = server.CreateServer()
	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/ws", websocketHandler)
	fmt.Println("Hosting server at port :" + port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

var wsServer *server.Server

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

func homeHandler(r http.ResponseWriter, req *http.Request) {
	http.ServeFile(r, req, "index.html")
}

func websocketHandler(r http.ResponseWriter, req *http.Request) {
	conn, err := upgrader.Upgrade(r, req, nil)
	if err != nil {
		return
	}
	wsServer.HandleNewConnection(conn)

}
