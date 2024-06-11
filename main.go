package main

import (
	"board/database"
	"board/handlers"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	database.InitDB()

	r := mux.NewRouter()
	r.HandleFunc("/threads", handlers.GetThreads).Methods("GET")
	r.HandleFunc("/threads", handlers.CreateThread).Methods("POST")
	r.HandleFunc("/replies", handlers.GetReplies).Methods("GET")
	r.HandleFunc("/replies", handlers.CreateReply).Methods("POST")

	log.Println("Server starting on :8080")
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatal(err)
	}
}
