package server

import (
	"time"
)

type Message struct {
	Id    string    `json:"id"`
	Email string    `json:"email"`
	Name  string    `json:"name"`
	Msg   string    `json:"msg"`
	Date  time.Time `json:"date"`
}

type Auth struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}
