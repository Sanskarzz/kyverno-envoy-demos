package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type User struct {
	FirstName string
	LastName  string
}

var users []User

func main() {

	router := mux.NewRouter
	users = append(users, User{FirstName: "Sanskar", LastName: "Gurdasani"})
	router.HandleFunc("/details", getuserdetails).Methods("GET")
	router.HandleFunc("/users", createuser).Methods("POST")

	fmt.Println("server starting on port 9000")
	log.Fatal(http.ListenAndServe(":9000", router))

}

func getuserdetails(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func createuser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var user User
	_ = json.NewDecoder(r.Body).Decode(&user)
	users = append(users, user)
	json.NewEncoder(w).Encode(&user)
}
