package main

import (
	//"fmt"
	"fmt"

	"github.com/mmcdole/gofeed"
)

// TODO
// 1. shift database to be an actual database
// 2. implement interface for multiple rss urls
// 3. implement a loop to check feeds every XX minutes
//
//

var database = make(map[string]*gofeed.Item)
var feedStore = []string{
	"http://www.commitstrip.com/en/feed/",
	"http://ryan.himmelwright.net/post/index.xml",
	"http://www.wuxiaworld.com/feed/"}

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
func storeFeed(feed *gofeed.Feed) bool {
	for i := 0; i < len(feed.Items); i++ {
		database[uniqueIdentifier(feed.Items[i])] = feed.Items[i]
	}
	return true
}

// Create a Feed and at it to DB
func addFeed(url string, feedparser *gofeed.Parser) bool {
	feed := createFeed(url, feedparser)
	return storeFeed(feed)
}

//iterate over all feed sources in feedStore
func addAllFeeds(feedparser *gofeed.Parser) bool {
	for _, element := range feedStore {
		addFeed(element, feedparser)
	}
	return true
}

func main() {
	feedparser := gofeed.NewParser()

	println(len(database))
	addAllFeeds(feedparser)
	//addFeed("http://www.wuxiaworld.com/feed/", feedparser)
	//addFeed("http://ryan.himmelwright.net/post/index.xml", feedparser)
	//addFeed("http://www.commitstrip.com/en/feed/", feedparser)
	println(len(database))
	println("hey its working.\n")
	fmt.Println("databse: ", database)
}
