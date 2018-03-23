package main

import (
	"github.com/mmcdole/gofeed"

	_ "github.com/mattn/go-sqlite3"
)

// TODO
// 1. [DONE] shift database to be an actual database
// 2. [DONE] implement interface for multiple rss urls
// 3. implement a loop to check feeds every XX minutes
//

func main() {
	debugPrint("Creating feed parser")
	feedparser := gofeed.NewParser()

	debugPrint("Initializing DB")
	db := initDB()

	debugPrint("Test Add to Feedstore")
	// FOR NOW, need this line when creating a new DB to add the test feeds
	testAddFeedSource(db)

	debugPrint("Adding feeds")
	addAllFeeds(feedparser, db)

	sqlTestPrint(db)

	debugPrint("hey its working.\n")
}
