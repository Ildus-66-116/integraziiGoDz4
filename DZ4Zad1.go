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

type MakeFriendsRequest struct {
	SourceID string `json:"source_id"`
	TargetID string `json:"target_id"`
}

func makeFriendsHandler(w http.ResponseWriter, r *http.Request) {
	var request MakeFriendsRequest

	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	sourceUser, sourceExists := users[request.SourceID]
	targetUser, targetExists := users[request.TargetID]

	if !sourceExists || !targetExists {
		http.Error(w, "One or both users do not exist", http.StatusNotFound)
		return
	}

	sourceUser.Friends = append(sourceUser.Friends, targetUser.Name)
	targetUser.Friends = append(targetUser.Friends, sourceUser.Name)

	users[request.SourceID] = sourceUser
	users[request.TargetID] = targetUser

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "%s и %s теперь друзья", sourceUser.Name, targetUser.Name)
}

func main() {
	http.HandleFunc("/create", createUserHandler)
	http.HandleFunc("/make_friends", makeFriendsHandler)

	fmt.Println("Starting server on :8080...")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err)
	}
}
