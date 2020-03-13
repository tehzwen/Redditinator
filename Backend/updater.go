package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"./db"
	"./reddit"
	"github.com/gocolly/colly"
)

func HandleUpdate(w http.ResponseWriter, r *http.Request, DB db.MyDB) {
	//update the posts first
	posts := UpdatePosts(DB)

	json.NewEncoder(w).Encode(posts)
}

func ScrapePage(post *reddit.SubredditPost, c *colly.Collector) {
	oldString := strings.Replace(post.FullLink, "www", "old", 1)

	c.OnHTML(".score", func(e *colly.HTMLElement) {
		if strings.Contains(e.Text, "upvoted") {
			fmt.Println("Here: ", e.Text)
			words := strings.Fields(e.Text)
			//fmt.Println(words[0])
			score, err := strconv.Atoi(words[0])
			if err != nil {
				panic(err)
			}
			post.Score = score
			fmt.Println(post.Score)

		}
	})

	c.Visit(oldString)
}

//Update the posts values for the past week
func UpdatePosts(DB db.MyDB) []reddit.SubredditPost {
	//get current time
	now := time.Now().Unix()
	then := now - 604800

	nowString := strconv.Itoa(int(now))
	thenString := strconv.Itoa(int(then))

	posts, err := DB.GetPostsBetweenDates(nowString, thenString)

	if err != nil {
		panic(err)
	}

	c := colly.NewCollector(
		colly.Async(false),
	)

	c.Limit(&colly.LimitRule{
		RandomDelay: 5 * time.Second,
	})

	c.OnError(func(r *colly.Response, err error) {
		fmt.Println(err)
	})

	//loop through posts and visit their main pages
	for i := 0; i < len(posts); i++ {
		ScrapePage(&posts[i], c)
	}
	//oldString := strings.Replace(posts[0].FullLink, "www", "old", 1)

	//fmt.Println(oldString)
	//ScrapePage(oldString, c)

	return posts
}
