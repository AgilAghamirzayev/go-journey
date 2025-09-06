package main

import (
	"encoding/json"
	"log"
	"net/http"
)

// User represents a user entity
type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

var users = []User{
	{ID: 1, Username: "user1", Email: "user1@example.com"},
	{ID: 2, Username: "user2", Email: "user2@example.com"},
}

func main() {
	http.HandleFunc("/users", getUsers)
	log.Println("Server started on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func getUsers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(users); err != nil {
		http.Error(w, "Failed to encode users", http.StatusInternalServerError)
		return
	}

}
