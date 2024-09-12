package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"strconv"
	"strings"

	"net/http"

	"github.com/gorilla/mux"
)

// The person Type
type Movie struct {
	ID        string `json:"id,omitempty"`
	Moviename string `json:"bookname,omitempty"`
	Actor     string `json:"author,omitempty"`
}

var movies []Movie

// Display all from the people var
func GetCollection(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(movies)
}

// Display a single data
func GetBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	for _, item := range movies {
		if item.ID == params["id"] {
			json.NewEncoder(w).Encode(item)
			return
		}
	}
	w.WriteHeader(http.StatusNotFound)
	json.NewEncoder(w).Encode(&Movie{})
}

// Create a new item
func CreateBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var movie Movie
	_ = json.NewDecoder(r.Body).Decode(&movie)
	movie.ID = strconv.Itoa(rand.Intn(1000000))
	movies = append(movies, movie)
	json.NewEncoder(w).Encode(&movie)
}

// Delete an item
func DeleteBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	for index, item := range movies {
		if item.ID == params["id"] {
			movies = append(movies[:index], movies[index+1:]...)
			break
		}
	}
	json.NewEncoder(w).Encode(movies)
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
	movies = append(movies, Movie{ID: "1", Moviename: "Inception", Actor: "Leonardo DiCaprio"})
	movies = append(movies, Movie{ID: "2", Moviename: "Batman", Actor: "Jack Nicholson"})
	router.HandleFunc("/", homePage)
	router.HandleFunc("/movie", GetCollection).Methods("GET")
	router.HandleFunc("/movie/{id}", GetBook).Methods("GET")
	router.HandleFunc("/movie", CreateBook).Methods("POST")
	router.HandleFunc("/movie/{id}", DeleteBook).Methods("DELETE")
	fmt.Println("server starting on port 8080")
	log.Fatal(http.ListenAndServe(":8080", router))

}
