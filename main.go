package main

import (
	"log"

	_ "github.com/mattn/go-sqlite3"
)

// TODO
// 1. [DONE] shift database to be an actual database
// 2. [DONE] implement interface for multiple rss urls
// 3. implement a loop to check feeds every XX minutes
//

var debug = true

// Prints only if debug global is true
func debugPrint(str string) {
	if debug {
		println(str)
	}
}

// Check error returns
func checkErr(err error) {
	if err != nil {
		log.Fatal("ERROR:", err)
	}
}

func main() {
	debugPrint("Initializing DB")
	db := initDB()

	addFeedSource("http://ryan.himmelwright.net/post/index.xml", "Test", db)

	updateAllFeedSources(db)
	debugPrint("hey its working.\n")

	startServer(db)
}
