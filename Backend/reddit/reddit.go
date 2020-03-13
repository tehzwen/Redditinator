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
	SentimentPos     float64 `db:"sentiment_pos"`
	SentimentNeg     float64 `db:"sentiment_neg"`
	SentimentNeu     float64 `db:"sentiment_neu"`
	SentimentOverall float64 `db:"sentiment_overall"`
}

type Reddit struct {
	TotalPosts  []SubredditPost
	CurrentTime int64
}

type Subreddit struct {
	ID   string `db:"id"`
	Name string `db:"name"`
}

type SubredditData struct {
	Data      []SubredditPost `json:"data"`
	Subreddit string
}

type CommentData struct {
	Data []PostComment `json:"data"`
}

type SubredditPost struct {
	Author        string `json:"author" db:"author"`
	TimeCreated   int64  `json:"created_utc" db:"created_utc"`
	FullLink      string `json:"full_link" db:"full_link"`
	ID            string `json:"id" db:"id"`
	IsVideo       bool   `json:"is_video" db:"is_video"`
	NumComments   int    `json:"num_comments" db:"num_comments"`
	NSFW          bool   `json:"over_18" db:"nsfw"`
	Score         int    `json:"score" db:"score"`
	SelfText      string `json:"selftext" db:"self_text"`
	SubredditID   string `json:"subreddit_id" db:"subreddit_id"`
	SubredditName string `json:"subreddit_name"`
	ThumbnailURL  string `json:"url" db:"thumbnail_url"`
	Title         string `json:"title" db:"title"`
	Sentiment     Sentiment
	Comments      []PostComment
}

//TODO add isPoster boolean to this struct

type PostComment struct {
	Author           string `json:"author" db:"author"`
	ID               string `json:"id" db:"id"`
	SubredditID      string `json:"subreddit_id" db:"subreddit_id"`
	Awards           int    `json:"total_awards_received"`
	Body             string `json:"body" db:"body"`
	PostID           string `json:"link_id" db:"post_id"`
	Score            int    `json:"score" db:"score"`
	Downs            int    `json:"downs" db:"downs"`
	IsPostAuthor     bool
	Controversiality int   `json:"controversiality" db:"controversy"`
	TimeCreated      int64 `json:"created_utc" db:"created_utc"`
	Sentiment        Sentiment
	ParentID         string `json:"parent_id" db:"parent_id"`
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
		r.TotalPosts[i].getCommentData("0")
	}
	wg.Done()
}

func (r *SubredditPost) getCommentData(score string) {
	if r.NumComments > 0 {
		request := "https://api.pushshift.io/reddit/comment/search/?link_id=" + r.ID + "&limit=1000"
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
				commentVal.Data[i].PostID = r.ID
				commentVal.Data[i].GetCommentSentiment()
			}
			r.Comments = append(r.Comments, commentVal.Data...)
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

	//iterate through and get sentiment for each post title
	for i := 0; i < len(postsData.Data); i++ {
		postsData.Data[i].GetPostSentiment()
		postsData.Data[i].SubredditName = subreddit
	}

	if len(postsData.Data) > 0 {
		r.TotalPosts = append(r.TotalPosts, postsData.Data...)
		r.CurrentTime = postsData.Data[len(postsData.Data)-1].TimeCreated
		return false
	}
	return true
}

func (d *SubredditPost) GetPostSentiment() {
	d.Sentiment = GetSentiment(d.Title)
}

func (c *PostComment) GetCommentSentiment() {
	c.Sentiment = GetSentiment(c.Body)
}

func GetSentiment(text string) Sentiment {
	//fmt.Println(text)
	analyzer := govader.NewSentimentIntensityAnalyzer()
	sentiment := analyzer.PolarityScores(text)
	return Sentiment{
		SentimentPos:     sentiment.Positive,
		SentimentNeg:     sentiment.Negative,
		SentimentNeu:     sentiment.Neutral,
		SentimentOverall: sentiment.Compound,
	}
}
