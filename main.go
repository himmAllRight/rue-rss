package main

import (
	"github.com/mmcdole/gofeed"
	//	"fmt"
)

var database []gofeed.Item

var testURL string = "http://www.wuxiaworld.com/feed/"

var feedparser = gofeed.NewParser()

func storeFeed(url string, destination []gofeed.Item) string {
	var feed, _ = feedparser.ParseURL(url)

	var i = 0 // counter
	for i = 0; i < len(feed.Items); i++ {
		destination.push(feed.Items[i])
	}
	return feed.Items[0].Title
}

func main() {

	print(len(database))
	var result = storeFeed(testURL, database)
	print(len(database))
	print("hey its working.\n")
	print(result)
	//getFeed(testURL)

}
