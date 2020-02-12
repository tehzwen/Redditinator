package main

import (
	"fmt"
	"net/http"

	"./db"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

func main() {

	r := mux.NewRouter()

	// Handle CORS
	headers := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"})
	methods := handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE"})
	origins := handlers.AllowedOrigins([]string{"*"})

	myDB := db.InitDB()

	/*
		Takes a post request body in the form of
		{
			"subreddits":["alberta", "witcher"],
			"before":"1234556",
			"after":"1231412123213"
		}
	*/
	r.HandleFunc("/subredditfetch", func(w http.ResponseWriter, r *http.Request) {
		CollectDataForSubreddits(w, r, myDB)
	}).Methods("POST")

	/*
		Takes a post request body in the form of
		{
			"text":"Example sentence here"
		}
	*/
	r.HandleFunc("/sentiment", func(w http.ResponseWriter, r *http.Request) {
		AnalyzeSentiment(w, r)
	}).Methods("POST")

	//WIP
	r.HandleFunc("/topic", func(w http.ResponseWriter, r *http.Request) {
		AnalyzeTopics(w, r)
	}).Methods("POST")

	//Queries - subreddit(optional search for posts of a certain subreddit)
	r.HandleFunc("/posts", func(w http.ResponseWriter, r *http.Request) {
		GetPosts(w, r, myDB)
	}).Methods("GET")

	/* Queries -
	subreddit(optional search for comments of a certain subreddit)
	postID(optional search for comments by postid)
	*/
	r.HandleFunc("/comments", func(w http.ResponseWriter, r *http.Request) {
		GetComments(w, r, myDB)
	}).Methods("GET")

	fmt.Println("Now serving on port 4000..")
	http.ListenAndServe(":4000", handlers.CORS(headers, methods, origins)(r))

	// //01/01/2017 https://www.unixtimestamp.com/index.php
	// //after := "1512086400"
	// //before := "1514764800"
}
