package handlers

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/mux"
	"go-server/models"
)

var (
	supabaseURL = os.Getenv("SUPABASE_URL")
	supabaseKey = os.Getenv("SUPABASE_KEY")
)

func GetTodos(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	req, _ := http.NewRequest("GET", supabaseURL+"/rest/v1/todos?select=*", nil)
	req.Header.Set("apikey", supabaseKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("GET /todos - Error fetching todos: %v", err)
		http.Error(w, "Failed to fetch todos", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	var todos []models.Todo
	json.NewDecoder(resp.Body).Decode(&todos)
	log.Printf("GET /todos - Returning %d todos", len(todos))
	json.NewEncoder(w).Encode(todos)
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

	body, _ := json.Marshal(todo)
	req, _ := http.NewRequest("POST", supabaseURL+"/rest/v1/todos", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("apikey", supabaseKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("POST /todos - Error creating todo: %v", err)
		http.Error(w, "Failed to create todo", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	var createdTodo models.Todo
	json.NewDecoder(resp.Body).Decode(&createdTodo)
	log.Printf("POST /todos - Created todo: ID=%d, Text=%s, Completed=%t", createdTodo.ID, createdTodo.Text, createdTodo.Completed)
	json.NewEncoder(w).Encode(createdTodo)
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

	var todo models.Todo
	err = json.NewDecoder(r.Body).Decode(&todo)
	if err != nil {
		log.Printf("PUT /todos/%d - Error decoding body: %v", id, err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	todo.ID = id

	body, _ := json.Marshal(todo)
	req, _ := http.NewRequest("PUT", supabaseURL+"/rest/v1/todos?id=eq."+strconv.Itoa(id), bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("apikey", supabaseKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("PUT /todos/%d - Error updating todo: %v", id, err)
		http.Error(w, "Failed to update todo", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	log.Printf("PUT /todos/%d - Updated todo: Text=%s, Completed=%t", id, todo.Text, todo.Completed)
	json.NewEncoder(w).Encode(todo)
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

	req, _ := http.NewRequest("GET", supabaseURL+"/rest/v1/todos?id=eq."+strconv.Itoa(id), nil)
	req.Header.Set("apikey", supabaseKey)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("PATCH /todos/%d/toggle - Error fetching todo: %v", id, err)
		http.Error(w, "Failed to fetch todo", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	var todos []models.Todo
	json.NewDecoder(resp.Body).Decode(&todos)
	if len(todos) == 0 {
		log.Printf("PATCH /todos/%d/toggle - Todo not found", id)
		http.Error(w, "Todo not found", http.StatusNotFound)
		return
	}

	todo := todos[0]
	todo.Completed = !todo.Completed

	body, _ := json.Marshal(todo)
	reqUpdate, _ := http.NewRequest("PUT", supabaseURL+"/rest/v1/todos?id=eq."+strconv.Itoa(id), bytes.NewBuffer(body))
	reqUpdate.Header.Set("Content-Type", "application/json")
	reqUpdate.Header.Set("apikey", supabaseKey)

	respUpdate, err := client.Do(reqUpdate)
	if err != nil {
		log.Printf("PATCH /todos/%d/toggle - Error updating todo: %v", id, err)
		http.Error(w, "Failed to toggle todo", http.StatusInternalServerError)
		return
	}
	defer respUpdate.Body.Close()

	log.Printf("PATCH /todos/%d/toggle - Toggled todo: Completed=%t", id, todo.Completed)
	json.NewEncoder(w).Encode(todo)
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

	req, _ := http.NewRequest("DELETE", supabaseURL+"/rest/v1/todos?id=eq."+strconv.Itoa(id), nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("apikey", supabaseKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("DELETE /todos/%d - Error deleting todo: %v", id, err)
		http.Error(w, "Failed to delete todo", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	log.Printf("DELETE /todos/%d - Deleted todo", id)
	w.WriteHeader(http.StatusNoContent)
}
