package main

type WebsocketMessage struct {
	Level string `json:"level"`
	Message string `json:"message"`
}