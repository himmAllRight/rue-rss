package main

import (
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"github.com/mmcdole/gofeed"
)

// DB Schema Def
var schema = `
CREATE TABLE IF NOT EXISTS person (
	first_name text,
	last_name text,
	email text
);

CREATE TABLE IF NOT EXISTS feedStore (
	feedurl TEXT, 
	category TEXT
);

CREATE TABLE IF NOT EXISTS feedData (
	feedname TEXT,
	feedurl TEXT,
	postname TEXT,
	posturl TEXT,
	publishdate TEXT,
	postdescription TEXT,
	postcontent TEXT
)`

// FeedSource struct that contains a feed source info
type FeedSource struct {
	Feedurl  string `db:"feedurl"`
	Category string `db:"category"`
}

// FeedItem struct that contains data for each feed item (ex: post)
type FeedItem struct {
	Feedname        string
	Feedurl         string
	Postname        string
	Posturl         string
	Publishdate     string
	Postdescription string
	Postcontent     string
}

// Init DB
func xinitDB() *sqlx.DB {
	db, err := sqlx.Connect("sqlite3", "test-db2.db")
	if err != nil {
		log.Fatalln(err)
	}

	db.MustExec(schema)

	return db
}

// Add a new feed source to the feedStore table
func xaddFeedSource(newURL string, category string, db *sqlx.DB) bool {
	feeds := []FeedSource{}
	db.Select(&feeds, "SELECT feedurl FROM feedStore where feedurl=$1", newURL)
	if len(feeds) == 0 {
		tx := db.MustBegin()
		tx.MustExec("INSERT INTO feedStore (feedurl, category) VALUES ($1, $2)", newURL, category)
		tx.Commit()
		return true
	}
	return false
}

// Create a feed object from a url string
func xcreateFeed(url string) (*gofeed.Feed, error) {
	feedparser := gofeed.NewParser()
	feed, err := feedparser.ParseURL(url)
	if err != nil {
		return nil, fmt.Errorf("Feed not found with parser")
	}
	return feed, nil
}

// Copies contents of Feed Item into the Feed Item struct
func createFeedItemStruct(feed *gofeed.Feed, feedItem *gofeed.Item) FeedItem {
	return FeedItem{feed.Title, feed.Link, feedItem.Title, feedItem.Link, feedItem.Published, feedItem.Description, feedItem.Content}
}

// Stores the feed item to the DB
// func storeFeedItem(feedItem FeedItem, db *sqlx.DB) bool {
// 	dbFeedItem := []FeedItem{}
// 	db.Select(&dbFeedItem, "SELECT posturl FROM feedData where posturl=$1", feedItem.Posturl)
// 	if len(dbFeedItem) == 0 {
// 		tx := db.MustBegin()
// 		tx.MustExec("INSERT INTO feedData (feedname, feedurl, postname, posturl, publishdate, postdescription, postcontent) VALUES ()", &feedItem)
// 		tx.Commit()
// 		return true
// 	}
// 	return false
// }

func sqlxTestMain() {
	db := xinitDB()

	xaddFeedSource("http://ryan.himmelwright.net/post/index.xml", "Test", db)

	feeds := []FeedSource{}
	err3 := db.Select(&feeds, "SELECT * FROM feedStore")

	fmt.Printf("Error: %+v\n", err3)
	fmt.Printf("Feeds: %+v\n", feeds)

	// feed, _ := xcreateFeed("http://ryan.himmelwright.net/post/index.xml")
	// feedItem := feed.Items[0]
	// feedObj := createFeedItemStruct(feed, feedItem)
	//storeFeedItem(feedObj, db)

}
