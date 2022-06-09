package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// creating our movie struct
type Movie struct {
	ID       string    `json:"id"`
	Isbn     string    `json:"isbn"`
	Title    string    `json:"title"`
	Director *Director `json:"director"`
}

// creating our director struct
type Director struct {
	Firstname string `json:"firstName"`
	Lastname  string `json:"lastName"`
}

// initializing our movie variable that will hold a slice of movies
var movies []Movie

func getMovies(w http.ResponseWriter, r *http.Request) {
	// setting content type as json
	w.Header().Set("Content-Type", "application/json")
	// encoding our response to json, our response is the movies slice
	json.NewEncoder(w).Encode(movies)
}

func getMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	params := mux.Vars(r)

	// looping through our slice to find the specific movie and sending it as an encoded response
	for _, item := range movies {
		if item.ID == params["id"] {
			json.NewEncoder(w).Encode(item)
			break
		}
	}
	json.NewEncoder(w).Encode(movies)
}

func deleteMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	// setting our params from the req params
	params := mux.Vars(r)

	// looping over our slice to find the movie with the id that matches the params
	for index, item := range movies {

		if item.ID == params["id"] {
			// basically removing the movie from the slice by appending the other movies and replacing its spot in the slice
			movies = append(movies[:index], movies[index+1:]...)
			break
		}
	}
}

func createMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	// creating a movie var that will hold our new movie
	var movie Movie
	// decoding our new movie from request and storing it in movie
	_ = json.NewDecoder(r.Body).Decode(&movie)
	// creating new id and converting it to string
	movie.ID = strconv.Itoa(rand.Intn(10000))
	// appending new movie to movies
	movies = append(movies, movie)
	// returning our newly made movie
	json.NewEncoder(w).Encode(movie)
}

func updateMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	params := mux.Vars(r)

	// looping over movies slice and deleting the movie that matches the params[id] then creating new movie from params. Wrong way to do it when dealing with databases but just a hack for this basic crud api
	for index, item := range movies {
		if item.ID == params["id"] {

			movies = append(movies[:index], movies[index+1:]...)
			var movie Movie
			_ = json.NewDecoder(r.Body).Decode(&movie)
			movie.ID = params["id"]
			movies = append(movies, movie)
			json.NewEncoder(w).Encode(movies)

			return
		}
	}

}

func main() {
	// creating a new instance of our mux router
	r := mux.NewRouter()

	// appending movies to our movie slice
	movies = append(movies, Movie{ID: "1", Isbn: "77234", Title: "Movie One", Director: &Director{Firstname: "John", Lastname: "Jacob"}}, Movie{ID: "2", Isbn: "57234", Title: "Movie Two", Director: &Director{Firstname: "Sam", Lastname: "Hunter"}})

	// route handling for each url path
	r.HandleFunc("/movies", getMovies).Methods("GET")
	r.HandleFunc("/movies/{id}", getMovie).Methods("GET")
	r.HandleFunc("/movies", createMovie).Methods("POST")
	r.HandleFunc("/movies/{id}", updateMovie).Methods("PUT")
	r.HandleFunc("/movies/{id}", deleteMovie).Methods("DELETE")

	fmt.Printf("Starting server at port 8000\n")
	log.Fatal(http.ListenAndServe(":8000", r))
}
