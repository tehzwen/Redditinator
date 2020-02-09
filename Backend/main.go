package main

import (
	"fmt"
	"net/http"

	"./db"
	"./reddit"
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

	//test adding a subreddit to the database
	myDB := db.InitDB()
	err := myDB.AddSubreddit(reddit.Subreddit{
		ID:   "12adw3",
		Name: "testsubreddit"})
	if err != nil {
		panic(err)
	}
	//test adding a post to the database
	err = myDB.AddPost(reddit.SubredditPost{
		Author:      "test",
		TimeCreated: 12341231,
		FullLink:    "https://test/com",
		ID:          "test_123675",
		IsVideo:     false,
		NumComments: 5,
		NSFW:        false,
		Score:       150,
		SelfText:    "Hey there buddy",
		SubredditID: "12adw3",
		Title:       "Look at my post!",
		Sentiment: reddit.Sentiment{
			SentimentPos:     1.0,
			SentimentNeg:     0.0,
			SentimentNeu:     0.0,
			SentimentOverall: 1.0,
		},
	})
	if err != nil {
		panic(err)
	}

	err = myDB.AddComment(reddit.PostComment{
		Author:           "Me",
		ID:               "890_cvbe",
		SubredditID:      "12adw3",
		Awards:           1,
		Body:             "this is a great post",
		PostID:           "test_123675",
		Score:            5,
		Downs:            0,
		Controversiality: 0,
		TimeCreated:      1235344324,
		Sentiment: reddit.Sentiment{
			SentimentPos:     1.0,
			SentimentNeg:     0.0,
			SentimentNeu:     0.0,
			SentimentOverall: 1.0,
		},
		IsPostAuthor: true,
	})

	if err != nil {
		panic(err)
	}

	http.ListenAndServe(":4000", handlers.CORS(headers, methods, origins)(r))

	// //01/01/2017 https://www.unixtimestamp.com/index.php
	// //after := "1512086400"
	// //before := "1514764800"
}
