package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"go-server/models"
)

func GetTodos(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	log.Printf("GET /todos - Returning %d todos", len(models.Todos))
	json.NewEncoder(w).Encode(models.Todos)
}

func CreateTodo(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var todo models.Todo
	err := json.NewDecoder(r.Body).Decode(&todo)
	if err != nil {
		log.Printf("POST /todos - Error decoding request body: %v", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	models.IDCounter++
	todo.ID = models.IDCounter
	models.Todos = append(models.Todos, todo)
	log.Printf("POST /todos - Created todo: ID=%d, Text=%s, Completed=%t", todo.ID, todo.Text, todo.Completed)
	json.NewEncoder(w).Encode(todo)
}

func UpdateTodo(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		log.Printf("PUT /todos/%s - Invalid ID: %v", params["id"], err)
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}
	for i, todo := range models.Todos {
		if todo.ID == id {
			err := json.NewDecoder(r.Body).Decode(&models.Todos[i])
			if err != nil {
				log.Printf("PUT /todos/%d - Error decoding request body: %v", id, err)
				http.Error(w, "Invalid request body", http.StatusBadRequest)
				return
			}
			models.Todos[i].ID = id
			log.Printf("PUT /todos/%d - Updated todo: Text=%s, Completed=%t", id, models.Todos[i].Text, models.Todos[i].Completed)
			json.NewEncoder(w).Encode(models.Todos[i])
			return
		}
	}
	log.Printf("PUT /todos/%d - Todo not found", id)
	http.Error(w, "Todo not found", http.StatusNotFound)
}

func ToggleTodo(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		log.Printf("PATCH /todos/%s/toggle - Invalid ID: %v", params["id"], err)
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}
	for i, todo := range models.Todos {
		if todo.ID == id {
			models.Todos[i].Completed = !models.Todos[i].Completed
			log.Printf("PATCH /todos/%d/toggle - Toggled todo: Completed=%t", id, models.Todos[i].Completed)
			json.NewEncoder(w).Encode(models.Todos[i])
			return
		}
	}
	log.Printf("PATCH /todos/%d/toggle - Todo not found", id)
	http.Error(w, "Todo not found", http.StatusNotFound)
}

func DeleteTodo(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		log.Printf("DELETE /todos/%s - Invalid ID: %v", params["id"], err)
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}
	for i, todo := range models.Todos {
		if todo.ID == id {
			models.Todos = append(models.Todos[:i], models.Todos[i+1:]...)
			log.Printf("DELETE /todos/%d - Deleted todo", id)
			w.WriteHeader(http.StatusNoContent)
			return
		}
	}
	log.Printf("DELETE /todos/%d - Todo not found", id)
	http.Error(w, "Todo not found", http.StatusNotFound)
}
