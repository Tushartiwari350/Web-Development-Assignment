package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/Tushar/shardBalancerAPI/HandlerFunc"
	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()
	fmt.Println("hello")
	r.HandleFunc("/shardbalancer/api", HandlerFunc.CheckHealth).Methods("GET")
	log.Fatal(http.ListenAndServe(":4000", r))
}
