package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os/exec"
	"runtime"
	"sync"

	"./db"
	"./reddit"
)

type CollectRequest struct {
	Subreddits []string `json:"subreddits"`
	From       int64    `json:"from"`
	To         int64    `json:"to"`
}

type AnalyzeRequest struct {
	Text  string `json:"text"`
	Topic string `json:"topic"`
}

func AnalyzeTopics(w http.ResponseWriter, r *http.Request) {

	anRequest := AnalyzeRequest{}
	err := json.NewDecoder(r.Body).Decode(&anRequest)
	if err != nil {
		fmt.Println(err)
	}
	const GOOS string = runtime.GOOS
	var pythonString string

	if GOOS == "linux" {
		pythonString = "python3"
	} else if GOOS == "windows" {
		pythonString = "py"
	}

	c := exec.Command(pythonString, "./LDA.py", anRequest.Text)
	out, err := c.Output()
	if err != nil {
		panic(err)
	}

	//output := string(out)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(string(out)))
}

func CollectDataForSubreddits(w http.ResponseWriter, r *http.Request, DB db.MyDB) {

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

	//go through the data and add it to the database correctly
	for x := 0; x < len(rClient.TotalPosts); x++ {
		err := DB.AddSubreddit(reddit.Subreddit{
			ID:   rClient.TotalPosts[x].SubredditID,
			Name: rClient.TotalPosts[x].SubredditName,
		})
		if err != nil {
			panic(err)
		}
		err = DB.AddPost(rClient.TotalPosts[x])
		if err != nil {
			panic(err)
		}

		for j := 0; j < len(rClient.TotalPosts[x].Comments); j++ {
			err = DB.AddComment(rClient.TotalPosts[x].Comments[j])
			if err != nil {
				panic(err)
			}
		}
	}
	fmt.Println("All done!")
	json.NewEncoder(w).Encode(rClient)
}

func GetPosts(w http.ResponseWriter, r *http.Request, DB db.MyDB) {
	params := r.URL.Query()
	subredditQuery := params.Get("subreddit")

	if subredditQuery == "" {
		vals, err := DB.GetPosts("")
		if err != nil {
			panic(err)
		}
		json.NewEncoder(w).Encode(vals)
	} else {
		vals, err := DB.GetPosts(subredditQuery)
		if err != nil {
			panic(err)
		}
		json.NewEncoder(w).Encode(vals)
	}
}

func GetComments(w http.ResponseWriter, r *http.Request, DB db.MyDB) {
	params := r.URL.Query()
	subredditQuery := params.Get("subreddit")
	postID := params.Get("postID")

	if subredditQuery != "" {
		vals, err := DB.GetComments(subredditQuery, "")
		if err != nil {
			panic(err)
		}
		json.NewEncoder(w).Encode(vals)
	} else if postID != "" {
		vals, err := DB.GetComments("", postID)
		if err != nil {
			panic(err)
		}
		json.NewEncoder(w).Encode(vals)
	} else {
		vals, err := DB.GetComments("", "")
		if err != nil {
			panic(err)
		}
		json.NewEncoder(w).Encode(vals)
	}
}

type SentimentRequest struct {
	Text string `json:"text"`
}

func AnalyzeSentiment(w http.ResponseWriter, r *http.Request) {
	req := SentimentRequest{}
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		fmt.Println(err)
	}
	sent := reddit.GetSentiment(req.Text)
	json.NewEncoder(w).Encode(sent)
}
