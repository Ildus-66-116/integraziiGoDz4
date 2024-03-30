package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"sync"
)

type User struct {
	ID      string   `json:"id"`
	Name    string   `json:"name"`
	Age     int      `json:"age"`
	Friends []string `json:"friends"`
}

var users = make(map[string]User)

var (
	userMutex = &sync.Mutex{}
	nextID    = 1
)

func createUserHandler(w http.ResponseWriter, r *http.Request) {
	var newUser User
	err := json.NewDecoder(r.Body).Decode(&newUser)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	userMutex.Lock()
	newUser.ID = strconv.Itoa(nextID)
	nextID++
	userMutex.Unlock()

	users[newUser.ID] = newUser

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"id": newUser.ID})
}

func main() {
	http.HandleFunc("/create", createUserHandler)

	fmt.Println("Starting server on :8080...")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err)
	}
	// curl.exe -d "{\"name\":\"Ildus\",\"age\":44,\"friends\":[\"3\"]}" -H "Content-Type: application/json" -X POST http://localhost:8080/create
	// curl.exe -d "{\"name\":\"Alena\",\"age\":42,\"friends\":[\"1\"]}" -H "Content-Type: application/json" -X POST http://localhost:8080/create
}
