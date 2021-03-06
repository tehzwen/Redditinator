package db

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"time"

	"../reddit"
	_ "github.com/denisenkom/go-mssqldb"
)

type DatabaseConnectionSecret struct {
	Password string `json:"password"`
	Server   string `json:"server"`
	User     string `json:"user"`
	Database string `json:"database"`
	Port     int
}

type MyDB struct {
	DB   *sql.DB
	Name string
}

func InitDB() MyDB {

	//first things first lets open the secret file
	jsonFile, err := os.Open("databaseSecret.json")
	if err != nil {
		panic(err)
	}

	defer jsonFile.Close()

	byteValue, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		panic(err)
	}

	var secret DatabaseConnectionSecret
	json.Unmarshal(byteValue, &secret)
	fmt.Println(secret)
	secret.Port = 1433

	connString := fmt.Sprintf("server=%s;user id=%s;password=%s;port=%d;database=%s", secret.Server, secret.User, secret.Password, secret.Port, secret.Database)

	conn, err := sql.Open("mssql", connString)
	if err != nil {
		log.Fatal("Open connection failed:", err.Error())
	}
	//defer conn.Close()

	return MyDB{
		DB:   conn,
		Name: "Topics",
	}
}

func (db *MyDB) AddSubreddit(s reddit.Subreddit) error {
	query := fmt.Sprintf(`addSubreddit '%s', '%s'`, s.ID, s.Name)
	_, err := db.DB.Exec(query)
	return err
}

func (db *MyDB) AddPost(p reddit.SubredditPost) error {
	titleString := strings.ReplaceAll(p.Title, "'", "")
	bodyString := strings.ReplaceAll(p.SelfText, "'", "")

	query := fmt.Sprintf(`addPost '%s', '%s', '%s', %d, '%s', %f, %f, %f, %f, %t, '%s', '%s', %d, '%s', %t, %d, '%s'`,
		p.ID, p.SubredditID, titleString, p.Score, p.Author, p.Sentiment.SentimentPos, p.Sentiment.SentimentNeg,
		p.Sentiment.SentimentNeu, p.Sentiment.SentimentOverall, p.NSFW, bodyString, p.ThumbnailURL, p.NumComments,
		p.FullLink, p.IsVideo, p.TimeCreated, p.Topic.String)

	_, err := db.DB.Exec(query)
	return err
}

func (db *MyDB) UpdateTopic(PostID string, topic string) error {
	query := fmt.Sprintf(`updateTopic '%s', '%s'`, PostID, topic)
	_, err := db.DB.Exec(query)
	return err
}

func (db *MyDB) UpdatePost(p reddit.SubredditPost) error {
	titleString := strings.ReplaceAll(p.Title, "'", "")
	bodyString := strings.ReplaceAll(p.SelfText, "'", "")

	query := fmt.Sprintf(`updatePost '%s', '%s', %d, '%s', %f, %f, %f, %f, %t, '%s', '%s', %d, '%s', %t, %d, '%s'`,
		p.ID, titleString, p.Score, p.Author, p.Sentiment.SentimentPos, p.Sentiment.SentimentNeg,
		p.Sentiment.SentimentNeu, p.Sentiment.SentimentOverall, p.NSFW, bodyString, p.ThumbnailURL, p.NumComments,
		p.FullLink, p.IsVideo, p.TimeCreated, p.Topic.String)

	fmt.Println(time.Now(), " updated post with id of ", p.ID)
	_, err := db.DB.Exec(query)
	return err
}

func (db *MyDB) AddComment(p reddit.PostComment) error {
	bodyString := strings.ReplaceAll(p.Body, "'", "")
	query := fmt.Sprintf(`addComment '%s', '%s', '%s', %d, '%s', %f, %f, %f, %f, '%s', %t, %d, %d, %d, %d, %s`,
		p.ID, p.PostID, p.SubredditID, p.Score, p.Author, p.Sentiment.SentimentPos, p.Sentiment.SentimentNeg, p.Sentiment.SentimentNeu,
		p.Sentiment.SentimentOverall, bodyString, p.IsPostAuthor, p.Awards, p.TimeCreated, p.Controversiality, p.Downs, p.ParentID)
	_, err := db.DB.Exec(query)
	return err
}

func (db *MyDB) GetPosts(subreddit string) ([]reddit.SubredditPost, error) {
	myPosts := []reddit.SubredditPost{}
	query := "SELECT * FROM post"
	if subreddit != "" {
		query = fmt.Sprintf("SELECT p.* FROM subreddit s, post p WHERE s.id = p.subreddit_id AND s.name ='%s'", subreddit)
	}
	rows, err := db.DB.Query(query)
	for rows.Next() {
		p := reddit.SubredditPost{}
		err := rows.Scan(&p.ID, &p.SubredditID, &p.Title, &p.Score, &p.Author, &p.Sentiment.SentimentPos, &p.Sentiment.SentimentNeg,
			&p.Sentiment.SentimentNeu, &p.Sentiment.SentimentOverall, &p.NSFW, &p.SelfText, &p.ThumbnailURL, &p.NumComments,
			&p.FullLink, &p.IsVideo, &p.TimeCreated, &p.Topic)

		if err != nil {
			panic(err)
		}
		myPosts = append(myPosts, p)
	}
	if err != nil {
		return nil, err
	}

	return myPosts, nil
}

func (db *MyDB) GetComments(subreddit, postID string) ([]reddit.PostComment, error) {
	myComments := []reddit.PostComment{}
	query := "SELECT * FROM comment"
	if postID != "" {
		query = fmt.Sprintf("SELECT c.* FROM post p, comment c WHERE p.id = c.post_id AND p.id ='%s'", postID)
	} else if subreddit != "" {
		query = fmt.Sprintf("SELECT c.* FROM subreddit s, comment c WHERE s.id = c.subreddit_id AND s.name ='%s'", subreddit)
	}
	rows, err := db.DB.Query(query)
	for rows.Next() {
		p := reddit.PostComment{}
		err := rows.Scan(&p.ID, &p.PostID, &p.SubredditID, &p.Score, &p.Author, &p.Sentiment.SentimentPos, &p.Sentiment.SentimentNeg, &p.Sentiment.SentimentNeu,
			&p.Sentiment.SentimentOverall, &p.Body, &p.IsPostAuthor, &p.Awards, &p.TimeCreated, &p.Controversiality, &p.Downs, &p.ParentID)

		if err != nil {
			panic(err)
		}
		myComments = append(myComments, p)
	}
	if err != nil {
		return nil, err
	}

	return myComments, nil
}

func (db *MyDB) GetSubredditNames() ([]string, error) {
	query := "SELECT name FROM subreddit"
	subreddits := []string{}

	rows, err := db.DB.Query(query)

	if err != nil {
		return subreddits, err
	}

	for rows.Next() {
		var temp string
		err := rows.Scan(&temp)

		if err != nil {
			return subreddits, err
		}
		subreddits = append(subreddits, temp)
	}

	return subreddits, nil
}

func (db *MyDB) GetTopLevelComments(postID string) ([]reddit.PostComment, error) {
	var query string
	if postID != "" {
		query = "SELECT * FROM comment WHERE parent_id LIKE 't3%' AND post_id='" + postID + "'"

	} else {
		query = "SELECT * FROM comment WHERE parent_id LIKE 't3%'"

	}
	comments := []reddit.PostComment{}
	rows, err := db.DB.Query(query)

	if err != nil {
		return comments, err
	}

	for rows.Next() {
		p := reddit.PostComment{}
		err := rows.Scan(&p.ID, &p.PostID, &p.SubredditID, &p.Score, &p.Author, &p.Sentiment.SentimentPos, &p.Sentiment.SentimentNeg, &p.Sentiment.SentimentNeu,
			&p.Sentiment.SentimentOverall, &p.Body, &p.IsPostAuthor, &p.Awards, &p.TimeCreated, &p.Controversiality, &p.Downs, &p.ParentID)
		if err != nil {
			return comments, err
		}
		comments = append(comments, p)
	}

	return comments, nil
}

func (db *MyDB) GetPostsBetweenDates(now, then string) ([]reddit.SubredditPost, error) {
	myPosts := []reddit.SubredditPost{}
	query := "SELECT * FROM post WHERE created_utc > " + then + " and created_utc < " + now

	rows, err := db.DB.Query(query)

	if err != nil {
		return myPosts, errors.New("Error running query: " + query)
	}

	for rows.Next() {
		p := reddit.SubredditPost{}
		err := rows.Scan(&p.ID, &p.SubredditID, &p.Title, &p.Score, &p.Author, &p.Sentiment.SentimentPos, &p.Sentiment.SentimentNeg,
			&p.Sentiment.SentimentNeu, &p.Sentiment.SentimentOverall, &p.NSFW, &p.SelfText, &p.ThumbnailURL, &p.NumComments,
			&p.FullLink, &p.IsVideo, &p.TimeCreated, &p.Topic)

		if err != nil {
			return myPosts, errors.New("Error scanning row")
		}
		myPosts = append(myPosts, p)
	}

	return myPosts, nil
}

type SubredditAvgSentiment struct {
	Name    string  `json:"name" db:"name"`
	Average float32 `json:"average" db:"average"`
}

func (db *MyDB) GetAverageSentimentOfSubreddits() ([]SubredditAvgSentiment, error) {
	mySentimentAverages := []SubredditAvgSentiment{}
	query := "SELECT s2.name, SUM(sentiment_overall)/COUNT(*) AS average FROM post, subreddit s2 WHERE post.subreddit_id =s2.id GROUP BY s2.name"
	rows, err := db.DB.Query(query)

	if err != nil {
		return mySentimentAverages, errors.New("Error running query: " + query)
	}

	for rows.Next() {
		p := SubredditAvgSentiment{}
		err := rows.Scan(&p.Name, &p.Average)

		if err != nil {
			return mySentimentAverages, errors.New("Error scanning row")
		}
		mySentimentAverages = append(mySentimentAverages, p)
	}

	return mySentimentAverages, nil
}

type AuthorCount struct {
	Author string `json:"author" db:"author"`
	Count  int    `json:"count" db:"cou"`
}

func (db *MyDB) GetTopAuthorsPerSubreddit(subredditName string) ([]AuthorCount, error) {

	myAuthorCounts := []AuthorCount{}
	query := "SELECT author, COUNT( author ) AS cou FROM post p, subreddit s WHERE p.subreddit_id=s.id AND s.name ='" + subredditName + "' GROUP BY author ORDER BY cou DESC"
	rows, err := db.DB.Query(query)

	if err != nil {
		return myAuthorCounts, errors.New("Error running query: " + query)
	}

	for rows.Next() {
		p := AuthorCount{}
		err := rows.Scan(&p.Author, &p.Count)

		if err != nil {
			return myAuthorCounts, errors.New("Error scanning row")
		}
		myAuthorCounts = append(myAuthorCounts, p)
	}

	return myAuthorCounts, nil
}

func (db *MyDB) SubredditSearch(searchValue string) ([]reddit.Subreddit, error) {
	mySubreddits := []reddit.Subreddit{}
	query := "SELECT * FROM subreddit WHERE name LIKE '%" + searchValue + "%' ORDER BY name"
	rows, err := db.DB.Query(query)

	if err != nil {
		fmt.Println(err)
		return mySubreddits, errors.New("Error running query: " + query)
	}

	for rows.Next() {
		r := reddit.Subreddit{}
		err := rows.Scan(&r.ID, &r.Name)

		if err != nil {
			return mySubreddits, errors.New("Error scanning row")
		}
		mySubreddits = append(mySubreddits, r)
	}

	return mySubreddits, nil
}

type TopicCount struct {
	Topic string `json:"topic" db:"topic"`
	Count int    `json:"count" db:"cou"`
}

func (db *MyDB) TopicOccurance(subredditID string) ([]TopicCount, error) {
	myTopics := []TopicCount{}
	query := "SELECT DISTINCT topic, COUNT(topic) AS cou FROM post WHERE post.subreddit_id='" + subredditID + "' AND topic != '' GROUP BY topic ORDER BY cou DESC"

	rows, err := db.DB.Query(query)

	if err != nil {
		fmt.Println(err)
		return myTopics, errors.New("Error running query: " + query)
	}

	for rows.Next() {
		t := TopicCount{}
		err := rows.Scan(&t.Topic, &t.Count)

		if err != nil {
			return myTopics, errors.New("Error scanning row")
		}
		myTopics = append(myTopics, t)
	}

	return myTopics, nil
}

type SubredditSentiment struct {
	PostSentiment    float64 `json:"postSent" db:"postSent"`
	PostPositive     float64 `json:"postSentPos" db:"postSentPos"`
	PostNegative     float64 `json:"postSentNeg" db:"postSentNeg"`
	PostNeutral      float64 `json:"postSentNeu" db:"postSentNeu"`
	CommentSentiment float64 `json:"commentSent" db:"commentSent"`
	CommentPositive  float64 `json:"commentSentPos" db:"commentSentPos"`
	CommentNegative  float64 `json:"commentSentNeg" db:"commentSentNeg"`
	CommentNeutral   float64 `json:"commentSentNeu" db:"commentSentNeu"`
}

func (db *MyDB) SubredditSentiment(subredditID string) (SubredditSentiment, error) {
	mySubredditSentiment := SubredditSentiment{}
	query := `SELECT * FROM (SELECT SUM(sentiment_overall)/COUNT(*) AS postSent, SUM(sentiment_pos )/COUNT(*) AS postSentPos, SUM(sentiment_neg )/COUNT(*) AS postSentNeg, SUM(sentiment_neu )/COUNT(*) AS postSentNeu FROM post p WHERE p.subreddit_id ='%s') AS p, 
(SELECT SUM(sentiment_overall)/COUNT(*) AS commentSent, SUM(sentiment_pos)/COUNT(*) AS commentSentPos, SUM(sentiment_neg)/COUNT(*) AS commentSentNeg, SUM(sentiment_neu)/COUNT(*) AS commentSentNeu FROM comment c WHERE c.subreddit_id='%s') AS c `

	query = fmt.Sprintf(query, subredditID, subredditID)
	rows, err := db.DB.Query(query)
	if err != nil {
		return mySubredditSentiment, errors.New("Error running query: " + query)
	}

	for rows.Next() {
		err := rows.Scan(&mySubredditSentiment.PostSentiment, &mySubredditSentiment.PostPositive, &mySubredditSentiment.PostNegative, &mySubredditSentiment.PostNeutral,
			&mySubredditSentiment.CommentSentiment, &mySubredditSentiment.CommentPositive, &mySubredditSentiment.CommentNegative, &mySubredditSentiment.CommentNeutral)
		if err != nil {
			fmt.Println(err)
			return mySubredditSentiment, errors.New("Error scanning row")
		}
	}

	return mySubredditSentiment, nil
}
