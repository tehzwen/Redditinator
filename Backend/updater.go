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

	// c.Limit(&colly.LimitRule{
	// 	//RandomDelay: time.Duration(rand.Intn(5)) * time.Second,
	// 	RandomDelay: 15 * time.Second,
	// 	Parallelism: 1,
	// })

	c.OnError(func(r *colly.Response, err error) {
		fmt.Println(err)
	})

	currentPost := &reddit.SubredditPost{}

	c.OnHTML("div[class=linkinfo]", func(e *colly.HTMLElement) {
		scoreText := e.ChildText(".number")
		val := strings.Replace(scoreText, ",", "", -1)
		score, err := strconv.Atoi(val)
		if err != nil {
			panic(err)
		}

		if score != 1 {
			currentPost.Score = score
			DB.UpdatePost((*currentPost))
		}
	})

	count := 0
	fmt.Println("Posts to scrape: ", len(posts))

	//loop through posts and visit their main pages
	for i := 0; i < len(posts); i++ {
		currentPost = &posts[i]
		oldString := strings.Replace(posts[i].FullLink, "www", "old", 1)
		c.Visit(oldString)
		count++

		progress := fmt.Sprintf("%.2f", (float32(count)/float32(len(posts)))*100)

		fmt.Println(progress, "% Complete")
	}
	c.Wait()

	return posts
}
