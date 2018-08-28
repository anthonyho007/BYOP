package main

import (
	"log"
	"net/http"
)

const (
	Port = "8000"
)

func main() {
	// http.HandleFunc("/", )
	log.Fatal(http.ListenAndServe(":"+Port, nil))
}
