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

func getFeedsForSource(feedSource FeedSource) {
	feedparser := gofeed.NewParser()
	feed, err := feedparser.ParseURL(feedSource.Feedurl)
	if err != nil {
		println("Error: %s", err)
	}
	//feedItems := []FeedItem{}
	for i := 0; i < len(feed.Items); i++ {
		println("Adding Feed Item: " + feed.Items[i].Title)
	}
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
// func createFeedItemStruct(feed *gofeed.Feed, feedItem *gofeed.Item) FeedItem {
// 	return FeedItem{feed.Title, feed.Link, feedItem.Title, feedItem.Link, feedItem.Published, feedItem.Description, feedItem.Content}
// }

// Iterates over all feed sources in feedStore table, and adds new feeds for each do the db
func updateAllFeedSources(db *sqlx.DB) {
	feedStore := []FeedSource{}
	db.Select(&feedStore, "SELECT * FROM feedStore")

	for _, feedSourceObj := range feedStore {
		debugPrint(feedSourceObj.Feedurl)
		feedSource, err := xcreateFeed(feedSourceObj.Feedurl)
		if err != nil {
			println("Error Creating Feed")
		}
		storeAllFeedItems(feedSource, db)
	}
}

// Stores all of the items for a feed source (if they don't exist)
func storeAllFeedItems(feed *gofeed.Feed, db *sqlx.DB) {
	// Iterate over feed items
	for i := 0; i < len(feed.Items); i++ {
		addedP := storeFeedItem(feed, feed.Items[i], db)
		if addedP {
			debugPrint("Feed Item Added: " + feed.Items[i].Title)
		} else {
			debugPrint("Feed Item Not Added: " + feed.Items[i].Title)
		}
	}
}

// Stores the feed item to the DB
func storeFeedItem(feed *gofeed.Feed, feedItem *gofeed.Item, db *sqlx.DB) bool {
	dbFeedItem := []FeedItem{}
	db.Select(&dbFeedItem, "SELECT posturl FROM feedData where posturl=$1", feedItem.Link)
	if len(dbFeedItem) == 0 {
		tx := db.MustBegin()
		tx.MustExec("INSERT INTO feedData (feedname, feedurl, postname, posturl, publishdate, postdescription, postcontent) VALUES (?, ?, ?, ?, ?, ?, ?)", feed.Title, feed.Link, feedItem.Title, feedItem.Link, feedItem.Published, feedItem.Description, feedItem.Content)
		tx.Commit()
		return true
	}
	return false
}

func sqlxTestMain() {
	db := xinitDB()

	xaddFeedSource("http://ryan.himmelwright.net/post/index.xml", "Test", db)

	updateAllFeedSources(db)

	storeAllFeedItems(feed, db)
}
