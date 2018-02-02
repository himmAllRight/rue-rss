package main

import (
	//"fmt"
	"github.com/mmcdole/gofeed"
)

var database = make(map[string]*gofeed.Item)
var testURL = "http://www.wuxiaworld.com/feed/"

var feedparser = gofeed.NewParser()

func uniqueIdentifier(feedItem *gofeed.Item) string {
	return feedItem.Link
}

func storeFeed(url string) string {
	var feed, _ = feedparser.ParseURL(url)

	var i = 0 // counter
	for i = 0; i < len(feed.Items); i++ {
		database[uniqueIdentifier(feed.Items[i])] = feed.Items[i]
	}
	return feed.Items[0].Title
}

func main() {

	print(len(database))
	var result = storeFeed(testURL)
	print(len(database))
	print("hey its working.\n")
	print(result)
	//getFeed(testURL)

}
