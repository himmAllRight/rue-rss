package main

import (
	"fmt"

	//"fmt"

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

var debug = true

// Test feed store to load up with new runs
var testFeedStore = []string{
	"http://www.commitstrip.com/en/feed/",
	"http://localhost:1313/post/index.xml",
	"http://ryan.himmelwright.net/post/index.xml",
	"http://www.wuxiaworld.com/feed/"}

// Prints only if debug global is true
func debugPrint(str string) {
	if debug {
		println(str)
	}
}

func initDB() *sql.DB {
	database, _ := sql.Open("sqlite3", "./testdb.db")

	// Create Feed Table
	feedStoreInitStatement, _ := database.Prepare("CREATE TABLE IF NOT EXISTS feedStore (id INTEGER PRIMARY KEY,feedurl TEXT, category TEXT)")
	feedStoreInitStatement.Exec()

	// Create Data Table
	statement, _ := database.Prepare("CREATE TABLE IF NOT EXISTS feedData (id INTEGER PRIMARY KEY,feedname TEXT,feedurl TEXT,postname TEXT,posturl TEXT,publishdate TEXT,postdescription TEXT,postcontent TEXT)")
	statement.Exec()

	// Return db ptr
	return (database)
}

// Test add Feed Source
func testAddFeedSource(db *sql.DB) bool {
	for _, url := range testFeedStore {
		addFeedSource(url, "cat", db)
	}
	return true
}

// takes a url that points to a feed and adds it to the the pool of feed sources
func addFeedSource(newURL string, category string, db *sql.DB) bool {
	// If the feed store doesn't contain feed, add it
	if !(sqlDoesContain("SELECT feedurl FROM feedStore WHERE feedurl='"+newURL+"'", db)) {
		if debug {
			println("Adding feed to feed store: ", newURL, " [", category, "]")
		}
		statement, _ := db.Prepare("INSERT INTO feedStore (feedurl, category) VALUES (?, ?)")
		statement.Exec(newURL, category)
	}
	return true
}

// Returns the store feeds
func getStoreFeeds(db *sql.DB) []string {
	debugPrint("Getting feeds from store")

	// Get feed urls from feedStore table
	var size int
	rows, _ := db.Query("SELECT feedurl FROM feedStore")
	defer rows.Close()
	numRows, _ := db.Query("SELECT COUNT(*) FROM feedStore")
	defer numRows.Close()
	numRows.Scan(&size)

	feedStore := make([]string, size)
	fmt.Println(size)
	// Add feeds to feedstore
	var feedurl string
	i := 0
	for rows.Next() {
		rows.Scan(&feedurl)
		fmt.Println("Feed url: " + feedurl)
		feedStore[i] = feedurl
		i = i + 1
	}

	return feedStore
}

// Create a new feed item from a url
func createFeed(url string, feedparser *gofeed.Parser) *gofeed.Feed {
	feed, _ := feedparser.ParseURL(url)
	return feed
}

// Add feed contents to DB
func storeFeed(url string, feed *gofeed.Feed, db *sql.DB) bool {
	feedPostURLs := getFeedPostURLs(url, db)
	if debug {
		println("\nStoring Feeds for: ", feed.Title)
		println("Feed URLs already in DB (So not adding):")
		printFeeds(feedPostURLs)
	}
	println("Adding Feeds:")
	for i := 0; i < len(feed.Items); i++ {
		feedItem := feed.Items[i]

		// If feed post not already in table, add it
		if !(sqlDoesContain("SELECT posturl FROM feedData WHERE posturl='"+feedItem.Link+"'", db)) {
			if debug {
				println("Adding entry: ", feedItem.Title, " [", feedItem.Link, "]")
			}
			statement, _ := db.Prepare("INSERT INTO feeddata (feedname, feedurl, postname, posturl, publishdate, postdescription, postcontent) VALUES (?, ?, ?, ?, ?, ?, ?)")
			statement.Exec(feed.Title, url, feedItem.Title, feedItem.Link, feedItem.Published, feedItem.Description, feedItem.Content)
		}
	}
	return true
}

// Create a Feed and at it to DB
func addFeed(url string, feedparser *gofeed.Parser, db *sql.DB) bool {
	feed := createFeed(url, feedparser)
	return storeFeed(url, feed, db)
}

//iterate over all feed sources in feedStore
func addAllFeeds(feedparser *gofeed.Parser, db *sql.DB) bool {
	// Get feeds from DB
	feedStore := getStoreFeeds(db)

	debugPrint("Feed Store")
	for _, element := range feedStore {
		debugPrint(element)
		fmt.Println("Feed url: " + element)
		addFeed(element, feedparser, db)
	}
	return true
}

// Gets all the post URLs for a feed
func getFeedPostURLs(feedURL string, db *sql.DB) *sql.Rows {
	rows, _ := db.Query("SELECT posturl FROM feeddata WHERE feedurl='" + feedURL + "'")
	return rows
}

func printFeeds(feedRows *sql.Rows) {
	var feedurl string
	for feedRows.Next() {
		feedRows.Scan(&feedurl)
		fmt.Println("\t" + feedurl)
	}
}

// Querys the DB and returns false if no results, or true if a result
func sqlDoesContain(query string, db *sql.DB) bool {
	rows, _ := db.Query(query)
	defer rows.Close()

	inResult := false

	for rows.Next() {
		inResult = true
	}

	return inResult
}

func sqlTestPrint(db *sql.DB) {
	rows, _ := db.Query("SELECT feedurl FROM feedStore")
	defer rows.Close()

	println("in sqlTestPrint")

	var feedurl string
	for rows.Next() {
		rows.Scan(&feedurl)
		fmt.Println("Feed url: " + feedurl)
	}
}

func main() {
	debugPrint("Creating feed parser")
	feedparser := gofeed.NewParser()

	debugPrint("Initializing DB")
	db := initDB()

	debugPrint("Test Add to Feedstore")
	testAddFeedSource(db)

	debugPrint("Adding feeds")
	addAllFeeds(feedparser, db)

	sqlTestPrint(db)

	debugPrint("hey its working.\n")
}
