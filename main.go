package main

import (
	"encoding/json"
	"net/http"
)

type Todo struct {
	ID        int    `json:"id"`
	Title     string `json:"title"`
	Completed bool   `json:"completed"`
}

var todos = []Todo{}
var nextID = 1

func main() {
	http.HandleFunc("/todos", todoHandler)
	http.ListenAndServe(":3008", nil)
}

func todoHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		json.NewEncoder(w).Encode(todos)
	case "POST":
		var t Todo
		json.NewDecoder(r.Body).Decode(&t)
		t.ID = nextID
		nextID++
		todos = append(todos, t)
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(t)
	}
}
