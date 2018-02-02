package main

import (
	//"fmt"
	"github.com/mmcdole/gofeed"
)

// TODO
// 1. shift database to be an actual database
// 2. implement interface for multiple rss urls
// 3. implement a loop to check feeds every XX minutes
//
//

var database = make(map[string]*gofeed.Item)
var testURL = "http://www.wuxiaworld.com/feed/"

var feedparser = gofeed.NewParser()

// generates the items key idetifier
func uniqueIdentifier(feedItem *gofeed.Item) string {
	return feedItem.Link
}

// Create a new feed item from a url
func createFeed(url string) *gofeed.Feed {
	feed, _ := feedparser.ParseURL(url)
	return feed
}

// Add feed contents to DB
func storeFeed(feed *gofeed.Feed) bool {
	for i := 0; i < len(feed.Items); i++ {
		database[uniqueIdentifier(feed.Items[i])] = feed.Items[i]
	}
	return true
}

// Create a Feed and at it to DB
func addFeed(url string) bool {
	feed := createFeed(url)
	return storeFeed(feed)
}

func main() {

	print(len(database))
	//var result = storeFeed(testURL)
	print(len(database))
	print("hey its working.\n")
	//print(result)
	//getFeed(testURL)

}
