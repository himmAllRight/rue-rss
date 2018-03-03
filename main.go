package main

import (
	//"fmt"
	"sync"

	"github.com/mmcdole/gofeed"

	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

// TODO
// 1. shift database to be an actual database
// 2. implement interface for multiple rss urls
// 3. implement a loop to check feeds every XX minutes
//
//

var database = make(map[string]*gofeed.Item)

var debug = true

var feedStore = []string{
	"http://www.commitstrip.com/en/feed/",
	"http://ryan.himmelwright.net/post/index.xml",
	"http://www.wuxiaworld.com/feed/"}

// Prints only if debug global is true
func debugPrint(str string) {
	if debug {
		println(str)
	}
}

func initDB() *sql.DB {
	debugPrint("Opening DB File")
	database, _ := sql.Open("sqlite3", "./testdb.db")

	// Create Table
	debugPrint("Creating Table")
	statement, _ := database.Prepare("CREATE TABLE IF NOT EXISTS feeddata (id ITNEGER PRIMARY KEY, feedname TEXT, url TEXT, read INTEGER)")
	debugPrint("Exec Table Creation")
	statement.Exec()

	// Return db ptr
	return (database)
}

// takes a url that points to a feed and adds it to the the pool of feed sources
func addFeedSource(newURL string) bool {
	feedStore = append(feedStore, newURL)
	return true
}

// generates the items key idetifier
func uniqueIdentifier(feedItem *gofeed.Item) string {
	return feedItem.Link
}

// Create a new feed item from a url
func createFeed(url string, feedparser *gofeed.Parser) *gofeed.Feed {
	feed, _ := feedparser.ParseURL(url)
	return feed
}

// Add feed contents to DB
func storeFeed(feed *gofeed.Feed, db *sql.DB) bool {
	if debug {
		println("Storing Feeds for: ", feed.Title)
	}
	for i := 0; i < len(feed.Items); i++ {
		if debug {
			println("Adding entry: ", feed.Items[i].Title)
		}

		// database[uniqueIdentifier(feed.Items[i])] = feed.Items[i]
		statement, _ := db.Prepare("INSERT INTO feeds (feedname, url, lastread) VALUES (?, ?, ?)")
		statement.Exec(feed.Title, feed.FeedLink, 0)
	}
	return true
}

// Create a Feed and at it to DB
func addFeed(url string, feedparser *gofeed.Parser, db *sql.DB) bool {
	feed := createFeed(url, feedparser)
	return storeFeed(feed, db)
}

//iterate over all feed sources in feedStore
func addAllFeeds(feedparser *gofeed.Parser, db *sql.DB) bool {
	var waitGroup sync.WaitGroup
	for _, element := range feedStore {
		waitGroup.Add(1)
		go func(element string, feedparser *gofeed.Parser) {
			addFeed(element, feedparser, db)
			waitGroup.Done()
		}(element, feedparser)
	}
	waitGroup.Wait()
	return true
}

func main() {
	debugPrint("Creating feed parser")
	feedparser := gofeed.NewParser()

	debugPrint("Initializing DB")
	db := initDB()

	debugPrint("Adding feeds")
	addAllFeeds(feedparser, db)

	debugPrint("hey its working.\n")
}
