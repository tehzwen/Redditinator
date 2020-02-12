package db

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

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

	query := fmt.Sprintf(`addPost '%s', '%s', '%s', %d, '%s', %f, %f, %f, %f, %t, '%s', '%s', %d, '%s', %t, %d`,
		p.ID, p.SubredditID, titleString, p.Score, p.Author, p.Sentiment.SentimentPos, p.Sentiment.SentimentNeg,
		p.Sentiment.SentimentNeu, p.Sentiment.SentimentOverall, p.NSFW, bodyString, p.ThumbnailURL, p.NumComments,
		p.FullLink, p.IsVideo, p.TimeCreated)
	_, err := db.DB.Exec(query)
	return err
}

func (db *MyDB) AddComment(p reddit.PostComment) error {
	bodyString := strings.ReplaceAll(p.Body, "'", "")
	query := fmt.Sprintf(`addComment '%s', '%s', '%s', %d, '%s', %f, %f, %f, %f, '%s', %t, %d, %d, %d, %d`,
		p.ID, p.PostID, p.SubredditID, p.Score, p.Author, p.Sentiment.SentimentPos, p.Sentiment.SentimentNeg, p.Sentiment.SentimentNeu,
		p.Sentiment.SentimentOverall, bodyString, p.IsPostAuthor, p.Awards, p.TimeCreated, p.Controversiality, p.Downs)
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
			&p.FullLink, &p.IsVideo, &p.TimeCreated)

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
			&p.Sentiment.SentimentOverall, &p.Body, &p.IsPostAuthor, &p.Awards, &p.TimeCreated, &p.Controversiality, &p.Downs)

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
