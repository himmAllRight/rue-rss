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
	statement, _ := database.Prepare("CREATE TABLE IF NOT EXISTS feeddata (id INTEGER PRIMARY KEY,feedname TEXT,feedurl TEXT,postname TEXT,posturl TEXT,publishdate TEXT,postdescription TEXT,postcontent TEXT)")
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
func storeFeed(url string, feed *gofeed.Feed, db *sql.DB) bool {
	feedPostURLs := getFeedPostURLs(url, db)
	if debug {
		println("Storing Feeds for: ", feed.Title)
		println("Feed URLs already in DB:")
		printFeeds(feedPostURLs)
	}
	for i := 0; i < len(feed.Items); i++ {
		feedItem := feed.Items[i]
		if debug {
			println("Adding entry: ", feedItem.Title)
		}

		// database[uniqueIdentifier(feed.Items[i])] = feed.Items[i]
		debugPrint("Prepare statement")
		statement, _ := db.Prepare("INSERT INTO feeddata (feedname, feedurl, postname, posturl, publishdate, postdescription, postcontent) VALUES (?, ?, ?, ?, ?, ?, ?)")
		debugPrint("Exec Insert")
		statement.Exec(feed.Title, url, feedItem.Title, feedItem.Link, feedItem.Published, feedItem.Description, feedItem.Content)
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
	for _, element := range feedStore {
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
		fmt.Println("Feed url: " + feedurl)
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
	rows, _ := db.Query("SELECT posturl FROM feeddata WHERE posturl='" + "http://ryan.himmelwright.net/post/solus-to-fedora/" + "'")
	defer rows.Close()

	var feedurl string
	for rows.Next() {
		rows.Scan(&feedurl)
		fmt.Println("Does this case hit?")
		fmt.Println("Feed url: " + feedurl)
	}
}

func main() {
	debugPrint("Creating feed parser")
	feedparser := gofeed.NewParser()

	debugPrint("Initializing DB")
	db := initDB()

	debugPrint("Adding feeds")
	addAllFeeds(feedparser, db)

	debugPrint("sqlTestPrint")
	sqlTestPrint(db)
	println(sqlDoesContain("SELECT posturl FROM feedData WHERE posturl='http://ryan.himmelwright.net/post/solus-to-fedora2/'", db))

	debugPrint("hey its working.\n")
}
