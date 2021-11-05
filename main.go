package main

import (
	// "fmt"
	"kv/handlers"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {

	// var store = map[string]interface{}{
	// 	"abc": 1,
	// 	"xyz": 2,
	// }

	// for k, v := range store {
	// 	fmt.Printf()
	// }

	r := mux.NewRouter()
	r.HandleFunc("/get", kv.GetAll).Methods("GET")
	r.HandleFunc("/get/{key}", kv.Get).Methods("GET")
	r.HandleFunc("/set", kv.Set).Methods("POST")
	r.HandleFunc("/search", kv.Search).Queries("prefix", "{str}").Methods("GET")
	r.HandleFunc("/search", kv.Search).Queries("suffix", "{str}").Methods("GET")

	log.Println("Listening on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", r))
}
