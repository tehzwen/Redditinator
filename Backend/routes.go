package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"

	"./reddit"
)

type CollectRequest struct {
	Subreddits []string `json:"subreddits"`
	From       int64    `json:"from"`
	To         int64    `json:"to"`
}

func CollectDataForSubreddits(w http.ResponseWriter, r *http.Request) {

	reqCollect := CollectRequest{}

	err := json.NewDecoder(r.Body).Decode(&reqCollect)

	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(reqCollect)

	after := "1580601600"
	before := "1580619600"

	//1580601600 //feb 2 2020
	//1580688000 //feb 3 2020
	//1580619600 //feb 2 2020 5am

	rClient := reddit.Reddit{}

	wg := sync.WaitGroup{}

	//loop through and get data for each of the subreddits
	for i := 0; i < len(reqCollect.Subreddits); i++ {
		go rClient.GetAllSubredditData("", before, after, reqCollect.Subreddits[i], &wg)
		wg.Add(1)
	}
	wg.Wait()

	fmt.Println("All done!")
	//print out the lengths of all the data retrieved by this reddit client
	// for key := range rClient.TotalPosts {
	// 	fmt.Printf("Total posts for %s: %d\n", key, len(rClient.TotalPosts[key]))
	// }

	json.NewEncoder(w).Encode(rClient)

}
