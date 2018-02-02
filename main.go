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

// stored all elements into the feed using their
func storeFeed(url string) bool {
	var feed, _ = feedparser.ParseURL(url)

	var i = 0 // counter
	for i = 0; i < len(feed.Items); i++ {
		database[uniqueIdentifier(feed.Items[i])] = feed.Items[i]
	}
	return true
}

func main() {

	print(len(database))
	//var result = storeFeed(testURL)
	print(len(database))
	print("hey its working.\n")
	//print(result)
	//getFeed(testURL)

}
