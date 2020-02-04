package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"

	"./reddit"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/jonreiter/govader"
)

func main() {

	r := mux.NewRouter()

	// Handle CORS
	headers := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"})
	methods := handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE"})
	origins := handlers.AllowedOrigins([]string{"*"})

	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {

		analyzer := govader.NewSentimentIntensityAnalyzer()
		sentiment := analyzer.PolarityScores("Alberta's energy sector leaves daunting environmental liabilities")
		fmt.Println("Compound score:", sentiment.Compound)
		fmt.Println("Positive score:", sentiment.Positive)
		fmt.Println("Neutral score:", sentiment.Neutral)
		fmt.Println("Negative score:", sentiment.Negative)

		after := "1580601600"
		before := "1580688000"

		rClient := reddit.Reddit{
			TotalPosts: make(map[string][]reddit.SubredditPost),
		}

		wg := sync.WaitGroup{}
		subreddits := []string{"alberta", "ProgrammerHumor", "witcher"}

		//loop through and get data for each of the subreddits
		for i := 0; i < len(subreddits); i++ {
			go rClient.GetAllSubredditData("", before, after, subreddits[i], &wg)
			wg.Add(1)
		}
		wg.Wait()

		fmt.Println("All done!")
		//print out the lengths of all the data retrieved by this reddit client
		for key := range rClient.TotalPosts {
			fmt.Printf("Total posts for %s: %d\n", key, len(rClient.TotalPosts[key]))
		}

		json.NewEncoder(w).Encode(rClient)

	}).Methods("GET")

	fmt.Println("Now serving on port 4000..")
	http.ListenAndServe(":4000", handlers.CORS(headers, methods, origins)(r))

	// //01/01/2017 https://www.unixtimestamp.com/index.php
	// //after := "1512086400"
	// //before := "1514764800"
}
