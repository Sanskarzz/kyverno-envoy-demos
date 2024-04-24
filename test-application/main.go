package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"math/rand"
	"strconv"
	"strings"

	"log"
	"net/http"

	"github.com/gorilla/mux"
)

// The person Type
type Book struct {
	ID       string `json:"id,omitempty"`
	Bookname string `json:"bookname,omitempty"`
	Author   string `json:"author,omitempty"`
}

var collection []Book

// Display all from the people var
func GetCollection(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(collection)
}

// Display a single data
func GetBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	for _, item := range collection {
		if item.ID == params["id"] {
			json.NewEncoder(w).Encode(item)
			return
		}
	}
	w.WriteHeader(http.StatusNotFound)
	json.NewEncoder(w).Encode(&Book{})
}

// Create a new item
func CreateBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var book Book
	_ = json.NewDecoder(r.Body).Decode(&book)
	book.ID = strconv.Itoa(rand.Intn(1000000))
	collection = append(collection, book)
	json.NewEncoder(w).Encode(&book)
}

// Delete an item
func DeleteBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	for index, item := range collection {
		if item.ID == params["id"] {
			collection = append(collection[:index], collection[index+1:]...)
			break
		}
	}
	json.NewEncoder(w).Encode(collection)
}

func homePage(w http.ResponseWriter, r *http.Request) {

	if basicAuth(w, r) {
		fmt.Fprintf(w, "Welcome to the HomePage!")
	} else {
		http.Error(w, "authorization failed", http.StatusUnauthorized)
		return
	}
}

func basicAuth(_ http.ResponseWriter, r *http.Request) bool {
	auth := strings.SplitN(r.Header.Get("Authorization"), " ", 2)

	if len(auth) != 2 || auth[0] != "Basic" {
		return false
	}

	payload, _ := base64.StdEncoding.DecodeString(auth[1])
	pair := strings.SplitN(string(payload), ":", 2)

	if len(pair) != 2 || !validate(pair[0], pair[1]) {
		return false
	}

	return true

}

func validate(username, password string) bool {
	if username == "test" && password == "test" {
		return true
	}
	return false
}

// main function to boot up everything
func main() {
	router := mux.NewRouter()
	collection = append(collection, Book{ID: "1", Bookname: "Harry Potter", Author: "J.K. Rowling"})
	collection = append(collection, Book{ID: "2", Bookname: "Animal Farm", Author: "George Orwell"})
	router.HandleFunc("/", homePage)
	router.HandleFunc("/book", GetCollection).Methods("GET")
	router.HandleFunc("/book/{id}", GetBook).Methods("GET")
	router.HandleFunc("/book", CreateBook).Methods("POST")
	router.HandleFunc("/book/{id}", DeleteBook).Methods("DELETE")
	log.Fatal(http.ListenAndServe(":8080", router))
}
