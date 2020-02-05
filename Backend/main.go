package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

func main() {

	r := mux.NewRouter()

	// Handle CORS
	headers := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"})
	methods := handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE"})
	origins := handlers.AllowedOrigins([]string{"*"})

	/*
		Takes a post request body in the form of
		{
			"subreddits":["alberta", "witcher"],
			"to":1234556,
			"from":1231412123213
		}
	*/
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		CollectDataForSubreddits(w, r)
	}).Methods("POST")

	fmt.Println("Now serving on port 4000..")
	http.ListenAndServe(":4000", handlers.CORS(headers, methods, origins)(r))

	// //01/01/2017 https://www.unixtimestamp.com/index.php
	// //after := "1512086400"
	// //before := "1514764800"
}
