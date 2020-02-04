package reddit

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"sync"
)

type Reddit struct {
	TotalPosts  map[string][]SubredditPost
	CurrentTime int64
}

type SubredditData struct {
	Data      []SubredditPost `json:"data"`
	Subreddit string
}

type SubredditPost struct {
	Author       string `json:"author"`
	TimeCreated  int64  `json:"created_utc"`
	FullLink     string `json:"full_link"`
	ID           string `json:"id"`
	IsVideo      bool   `json:"is_video"`
	NumComments  int    `json:"num_comments"`
	NSFW         bool   `json:"over_18"`
	Score        int    `json:"score"`
	SelfText     string `json:"selftext"`
	Subreddit    string `json:"subreddit"`
	ThumbnailURL string `json:"url"`
	Title        string `json:"title"`
}

func (r *Reddit) GetAllSubredditData(query, before, after, subreddit string, wg *sync.WaitGroup) {
	r.getSubredditData("", before, after, subreddit)

	for {
		done := r.getSubredditData("", before, strconv.FormatInt(r.CurrentTime, 10), subreddit)
		if done {
			wg.Done()
			break
		}
	}
}

func (r *Reddit) getSubredditData(query, before, after, subreddit string) bool {
	request := "https://api.pushshift.io/reddit/search/submission/?q=" + query +
		"&after=" + after + "&before=" + before + "&subreddit=" + subreddit + "&size=1000"

	fmt.Println("Fetching...", request)
	res, err := http.Get(request)

	if err != nil {
		fmt.Println("Error: ", err)
	}

	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)

	if err != nil {
		fmt.Println("Error: ", err)
	}

	postsData := SubredditData{
		Subreddit: subreddit,
	}

	err = json.Unmarshal(body, &postsData)

	if err != nil {
		fmt.Println("Error: ", err)
	}

	if len(postsData.Data) <= 0 {
		return true
	}

	r.TotalPosts[subreddit] = append(r.TotalPosts[subreddit], postsData.Data...)
	r.CurrentTime = postsData.Data[len(postsData.Data)-1].TimeCreated
	return false
}
