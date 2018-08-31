package server

import (
	"time"
)

type Message struct {
	Id   string    `json:"id"`
	Code string    `json:"code"`
	Name string    `json:"name"`
	Msg  string    `json:"msg"`
	Date time.Time `json:"date"`
}

type Auth struct {
	Name string `json:"name"`
	Code string `json:"code"`
}
