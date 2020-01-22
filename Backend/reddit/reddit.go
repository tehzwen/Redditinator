package reddit

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
)

type Reddit struct {
	Auth        AccessToken
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

// NOT IN USE { REFERENCE }
//struct to hold json info from file that contains keys
type RedditKeys struct {
	Username     string `json:"username"`
	Password     string `json:"password"`
	AppID        string `json:"appID"`
	AppSecret    string `json:"appSecret"`
	AppUserAgent string `json:"appUserAgent"`
}

// NOT IN USE { REFERENCE }
type AccessToken struct {
	Token      string `json:"access_token"`
	Expiration int    `json:"expires_in"`
	Scope      string `json:"scope"`
	TokenType  string `json:"token_type"`
	UserAgent  string
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

// NOT IN USE { REFERENCE }
func refreshAccessToken(r *Reddit) {
	//read in from keys json file
	keyFile, err := os.Open("./reddit/keys.json")
	if err != nil {
		fmt.Println("Error: ", err)
	}

	byteVal, err := ioutil.ReadAll(keyFile)
	if err != nil {
		fmt.Println("Error: ", err)
	}

	var keyData RedditKeys
	json.Unmarshal(byteVal, &keyData)

	if keyData.AppID != "" {
		bodyData := strings.NewReader("grant_type=password&username=" + keyData.Username + "&password=" + keyData.Password)
		req, err := http.NewRequest("POST", "https://www.reddit.com/api/v1/access_token", bodyData)
		if err != nil {
			fmt.Println("Error: ", err)
		}
		req.SetBasicAuth(keyData.AppID, keyData.AppSecret)
		req.Header.Set("User-Agent", keyData.AppUserAgent)

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			fmt.Println("Error: ", err)
		}

		if resp.StatusCode == 200 {
			var t AccessToken

			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				panic(err.Error())
			}

			err = json.Unmarshal(body, &t)

			if err != nil {
				panic(err.Error())
			}

			fmt.Println(t, resp.StatusCode)
			t.UserAgent = keyData.AppUserAgent

			r.Auth = t
		}
		defer resp.Body.Close()
	}

}
