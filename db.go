package main

import (
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"github.com/mmcdole/gofeed"
)

// DB Schema Def
var schema = `
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

// Init DB// Updates all the feed sources in feedStore table
func initDB() *sqlx.DB {
	db, err := sqlx.Connect("sqlite3", "test-db2.db")
	checkErr(err)
	db.MustExec(schema)

	return db
}

// Add a new feed source to the feedStore table
func addFeedSource(newURL string, category string, db *sqlx.DB) bool {
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

// Removes a source feed from the feedStore table, and it's associated data
func deleteFeedSource(feedurl string, db *sqlx.DB) {
	tx := db.MustBegin()
	tx.MustExec("DELETE FROM feedStore WHERE feedurl=?;", feedurl)
	tx.MustExec("DELETE FROM feedData WHERE feedurl=?;", feedurl)
	tx.Commit()
}

// Create a feed object from a url string
func createFeed(url string) (*gofeed.Feed, error) {
	feedparser := gofeed.NewParser()
	feed, err := feedparser.ParseURL(url)
	checkErr(err)
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
		storeAllFeedItems(feedSourceObj, db)
	}
}

// Stores all of the items for a feed source (if they don't exist)
func storeAllFeedItems(feedSource FeedSource, db *sqlx.DB) {
	feed, err := createFeed(feedSource.Feedurl)
	checkErr(err)

	// Iterate over feed items
	for i := 0; i < len(feed.Items); i++ {
		addedP := storeFeedItem(feedSource, feed, feed.Items[i], db)
		if addedP {
			debugPrint("Feed Item Added: " + feed.Items[i].Title)
		} else {
			debugPrint("Feed Item Not Added: " + feed.Items[i].Title)
		}
	}
}

// Returns a feedItem object, if it exists in the DB (feedData Table)
func getFeedItemData(posturl string, db *sqlx.DB) (FeedItem, error) {
	dbFeedItem := FeedItem{}
	db.Get(&dbFeedItem, "SELECT * FROM feedData where posturl=$1", posturl)

	fmt.Printf("posturl: %s\n", posturl)
	fmt.Printf("return struct:\n%+v\n", dbFeedItem)

	if dbFeedItem.Posturl != "" {
		return dbFeedItem, nil
	}
	return dbFeedItem, errors.New("No match ")
}

// Stores the feed item to the DB
func storeFeedItem(feedSource FeedSource, feed *gofeed.Feed, feedItem *gofeed.Item, db *sqlx.DB) bool {
	dbFeedItem := []FeedItem{}
	db.Select(&dbFeedItem, "SELECT posturl FROM feedData where posturl=$1", feedItem.Link)
	if len(dbFeedItem) == 0 {
		tx := db.MustBegin()
		tx.MustExec("INSERT INTO feedData (feedname, feedurl, postname, posturl, publishdate, postdescription, postcontent) VALUES (?, ?, ?, ?, ?, ?, ?)", feed.Title, feedSource.Feedurl, feedItem.Title, feedItem.Link, feedItem.Published, feedItem.Description, feedItem.Content)
		tx.Commit()
		return true
	}
	return false
}

func sqlxTestMain() {
	db := initDB()

	addFeedSource("http://ryan.himmelwright.net/post/index.xml", "Test", db)

	updateAllFeedSources(db)

}
