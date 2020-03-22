package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"

	"./db"
	"./reddit"
	"github.com/gorilla/mux"
)

type CollectRequest struct {
	Subreddits []string `json:"subreddits"`
	Before     string   `json:"before"`
	After      string   `json:"after"`
}

type AnalyzeRequest struct {
	PostID string `json:"id"`
	Topic  string `json:"topic"`
}

func AnalyzeTopics(w http.ResponseWriter, r *http.Request, DB db.MyDB) {
	anRequest := AnalyzeRequest{}
	err := json.NewDecoder(r.Body).Decode(&anRequest)
	if err != nil {
		fmt.Println(err)
	}

	if anRequest.PostID == "" || anRequest.Topic == "" {
		w.Write([]byte("Need to provide body fields"))

	} else {
		err := DB.UpdateTopic(anRequest.PostID, anRequest.Topic)
		if err != nil {
			panic(err)
		}
	}

}

func CollectDataForSubreddits(w http.ResponseWriter, r *http.Request, DB db.MyDB) {

	reqCollect := CollectRequest{}

	err := json.NewDecoder(r.Body).Decode(&reqCollect)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(reqCollect)

	// after := "1580601600"
	// before := "1580619600"

	after := reqCollect.After
	before := reqCollect.Before

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

func GetAuthors(w http.ResponseWriter, r *http.Request, DB db.MyDB) {
	params := r.URL.Query()
	subredditQuery := params.Get("subreddit")

	if subredditQuery == "" {
		w.Write([]byte("Need to provide subreddit query field"))
	} else {
		vals, err := DB.GetTopAuthorsPerSubreddit(subredditQuery)
		if err != nil {
			panic(err)
		}
		json.NewEncoder(w).Encode(vals)
	}
}

func GetSubreddits(w http.ResponseWriter, r *http.Request, DB db.MyDB) {
	params := r.URL.Query()
	subredditQuery := params.Get("subreddit")

	if subredditQuery == "" {
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		searched, err := DB.SubredditSearch(subredditQuery)
		if err != nil {
			fmt.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
		} else {
			json.NewEncoder(w).Encode(searched)
		}
	}
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
	topLevel := params.Get("topLevel")

	if subredditQuery != "" {
		vals, err := DB.GetComments(subredditQuery, "")
		if err != nil {
			panic(err)
		}
		json.NewEncoder(w).Encode(vals)
	} else if postID != "" {
		if topLevel != "" {
			vals, err := DB.GetTopLevelComments(postID)
			if err != nil {
				panic(err)
			}
			json.NewEncoder(w).Encode(vals)
		} else {
			vals, err := DB.GetComments("", postID)
			if err != nil {
				panic(err)
			}
			json.NewEncoder(w).Encode(vals)
		}

	} else {
		if topLevel != "" {
			vals, err := DB.GetTopLevelComments("")
			if err != nil {
				panic(err)
			}
			json.NewEncoder(w).Encode(vals)
			LogComments(vals)
		} else {

			vals, err := DB.GetComments("", "")
			if err != nil {
				panic(err)
			}
			json.NewEncoder(w).Encode(vals)
			LogComments(vals)
		}
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

func UpdatePost(w http.ResponseWriter, r *http.Request, DB db.MyDB) {
	req := reddit.SubredditPost{}
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		fmt.Println(err)
	}
	err = DB.UpdatePost(req)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Println(err)
	}
	w.WriteHeader(http.StatusOK)
}

func DailyRedditFetch(w http.ResponseWriter, r *http.Request, DB db.MyDB) {
	subs, err := DB.GetSubredditNames()
	if err != nil {
		panic(err)
	}

	rClient := reddit.Reddit{}
	//get the current time
	currTime := time.Now().Unix()
	before := strconv.FormatInt(currTime, 10)
	after := strconv.FormatInt(currTime-86400, 10)
	wg := sync.WaitGroup{}

	//loop through and get data for each of the subreddits
	for i := 0; i < len(subs); i++ {
		go rClient.GetAllSubredditData("", before, after, subs[i], &wg)
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
	lengthString := strconv.Itoa(len(rClient.TotalPosts))
	testData := []byte("Daily fetch of data completed successfully. Gathered " + lengthString + " posts.")
	f, err := os.OpenFile("status.log",
		os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Println(err)
	}
	defer f.Close()

	logger := log.New(f, "Reddit Fetch Log: ", log.LstdFlags)
	logger.Println(string(testData))

	if err != nil {
		panic(err)
	}

	json.NewEncoder(w).Encode(rClient)
}

func LogComments(vals []reddit.PostComment) {
	lengthString := strconv.Itoa(len(vals))
	testData := []byte("Retrieved " + lengthString + " comments.")
	f, err := os.OpenFile("status.log",
		os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Println(err)
	}
	defer f.Close()

	logger := log.New(f, "Comment Fetch: ", log.LstdFlags)
	logger.Println(string(testData))
}

func GetTopicOccurance(w http.ResponseWriter, r *http.Request, DB db.MyDB) {
	vars := mux.Vars(r)
	occurance, err := DB.TopicOccurance(vars["id"])
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		json.NewEncoder(w).Encode(occurance)
	}
}

func GetSubredditSentiment(w http.ResponseWriter, r *http.Request, DB db.MyDB) {
	vars := mux.Vars(r)
	occurance, err := DB.SubredditSentiment(vars["id"])
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		json.NewEncoder(w).Encode(occurance)
	}
}
