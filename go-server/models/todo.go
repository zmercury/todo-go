package models

import "log"

type Todo struct {
	ID        int    `json:"id"`
	Text      string `json:"text"`
	Completed bool   `json:"completed"`
}

var Todos []Todo

var IDCounter int

func InitTodos() {
	Todos = []Todo{
		{
			ID:        1,
			Text:      "Start with this todo",
			Completed: false,
		},
	}
	IDCounter = 1
	log.Printf("Initialized with todo: ID=%d, Text=%s, Completed=%t", Todos[0].ID, Todos[0].Text, Todos[0].Completed)
}
