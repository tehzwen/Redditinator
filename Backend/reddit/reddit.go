package reddit

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"sync"

	"github.com/jonreiter/govader"
)

type Sentiment struct {
	SentimentPos     float64
	SentimentNeg     float64
	SentimentNeu     float64
	SentimentOverall float64
}

type Reddit struct {
	TotalPosts  []SubredditPost
	CurrentTime int64
}

type SubredditData struct {
	Data      []SubredditPost `json:"data"`
	Subreddit string
}

type CommentData struct {
	Data []PostComment `json:"data"`
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
	SubredditID  string `json:"subreddit_id"`
	ThumbnailURL string `json:"url"`
	Title        string `json:"title"`
	Sentiment    Sentiment
	Comments     []PostComment
}

type PostComment struct {
	Author           string `json:"author"`
	ID               string `json:"id"`
	SubredditID      string `json:"subreddit_id"`
	Awards           int    `json:"total_awards_received"`
	Body             string `json:"body"`
	PostID           string `json:"link_id"`
	Score            int    `json:"score"`
	Downs            int    `json:"downs"`
	Controversiality int    `json:"controversiality"`
	TimeCreated      int64  `json:"created_utc"`
	Sentiment        Sentiment
}

func (r *Reddit) GetAllSubredditData(query, before, after, subreddit string, wg *sync.WaitGroup) {
	r.getSubredditData("", before, after, subreddit)
	wg.Add(1)

	for {
		done := r.getSubredditData("", before, strconv.FormatInt(r.CurrentTime, 10), subreddit)
		if done {
			wg.Done()
			break
		}
	}

	//once done we need to go analyze the comments for each of the post
	for i := 0; i < len(r.TotalPosts); i++ {
		r.TotalPosts[i].getCommentData("1")
	}
	wg.Done()
}

func (r *SubredditPost) getCommentData(score string) {
	request := "https://api.pushshift.io/reddit/comment/search/?link_id=" + r.ID + "&limit=1000&sort=desc&sort_type=score&score=%3E" + score
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

	commentVal := CommentData{}

	err = json.Unmarshal(body, &commentVal)

	if err != nil {
		fmt.Println("ERROR PARSING ", commentVal)
	} else if len(commentVal.Data) > 0 {
		//perform sentiment on the comments
		for i := 0; i < len(commentVal.Data); i++ {
			commentVal.Data[i].GetCommentSentiment()
		}
		r.Comments = append(r.Comments, commentVal.Data...)
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

	//iterate through and get sentiment for each post title
	for i := 0; i < len(postsData.Data); i++ {
		postsData.Data[i].GetPostSentiment()
	}

	if len(postsData.Data) > 0 {
		r.TotalPosts = append(r.TotalPosts, postsData.Data...)
		r.CurrentTime = postsData.Data[len(postsData.Data)-1].TimeCreated
		return false
	}
	return true
}

func (d *SubredditPost) GetPostSentiment() {
	analyzer := govader.NewSentimentIntensityAnalyzer()
	sentiment := analyzer.PolarityScores(d.Title)
	d.Sentiment.SentimentPos = sentiment.Positive
	d.Sentiment.SentimentNeg = sentiment.Negative
	d.Sentiment.SentimentNeu = sentiment.Neutral
	d.Sentiment.SentimentOverall = sentiment.Compound
}

func (c *PostComment) GetCommentSentiment() {
	analyzer := govader.NewSentimentIntensityAnalyzer()
	sentiment := analyzer.PolarityScores(c.Body)
	c.Sentiment.SentimentPos = sentiment.Positive
	c.Sentiment.SentimentNeg = sentiment.Negative
	c.Sentiment.SentimentNeu = sentiment.Neutral
	c.Sentiment.SentimentOverall = sentiment.Compound
}
