package main

import (
	"log"
	"net/http"

	"github.com/Tushar/myapi/HandlerFunc"
	"github.com/gorilla/mux"
)

func Newrouter() {
	r := mux.NewRouter()
	r.HandleFunc("/twitter/users", HandlerFunc.AddPost).Methods("POST")
	r.HandleFunc("/twitter/users/{name}", HandlerFunc.GetMyPosts).Methods("GET")
	r.HandleFunc("/twitter/follow", HandlerFunc.Follow).Methods("POST")
	r.HandleFunc("/twitter/read/{name}", HandlerFunc.ReadPost).Methods("GET")

	log.Fatal(http.ListenAndServe(":8080", r))
}
