package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"go-server/handlers"
	"go-server/middleware"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	router := mux.NewRouter()
	router.Use(middleware.LogRequest)

	router.HandleFunc("/todos", handlers.GetTodos).Methods("GET")
	router.HandleFunc("/todos", handlers.CreateTodo).Methods("POST")
	router.HandleFunc("/todos/{id}", handlers.UpdateTodo).Methods("PUT")
	router.HandleFunc("/todos/{id}/toggle", handlers.ToggleTodo).Methods("PATCH")
	router.HandleFunc("/todos/{id}", handlers.DeleteTodo).Methods("DELETE")

	log.Println("Server starting on http://localhost:6969")
	err = http.ListenAndServe(":6969", router)
	if err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}

