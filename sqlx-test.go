package main

import (
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

// DB Schema Def
var schema = `
CREATE TABLE IF NOT EXISTS person (
	first_name text,
	last_name text,
	email text
);

CREATE TABLE IF NOT EXISTS feedStore (
	id INTEGER PRIMARY KEY,
	feedurl TEXT, 
	category TEXT
);

CREATE TABLE IF NOT EXISTS feedData (
	id INTEGER PRIMARY KEY,
	feedname TEXT,
	feedurl TEXT,
	postname TEXT,
	posturl TEXT,
	publishdate TEXT,
	postdescription TEXT,
	postcontent TEXT
)`

// FeedSource struct that contains a feed source info
type FeedSource struct {
	ID       int    `db:"id"`
	Feedurl  string `db:"feedurl"`
	Category string `db:"category"`
}

// FeedItem struct that contains data for each feed item (ex: post)
type FeedItem struct {
	ID              int
	Feedname        string
	Feedurl         string
	Postname        string
	Posturl         string
	Publishdate     string
	Postdescription string
	Postcontent     string
}

// Init DB
func xinitDB() (*sqlx.DB, *sqlx.Tx) {
	db, err := sqlx.Connect("sqlite3", "test-db2.db")
	if err != nil {
		log.Fatalln(err)
	}

	db.MustExec(schema)

	tx := db.MustBegin()
	return db, tx
}

// Add a new feed source to the feedStore table
func xaddFeedSource(newURL string, category string, db *sqlx.DB, tx *sqlx.Tx) bool {
	feeds := []FeedSource{}
	db.Select(&feeds, "SELECT feedurl FROM feedStore where feedurl=$1", newURL)
	if len(feeds) == 0 {
		tx.MustExec("INSERT INTO feedStore (feedurl, category) VALUES ($1, $2)", newURL, category)
		tx.Commit()
		return true
	}
	return false
}

func sqlxTestMain() {
	db, tx := xinitDB()

	xaddFeedSource("http://ryan.himmelwright.net/post/index.xml", "Test", db, tx)

	feeds := []FeedSource{}
	err3 := db.Select(&feeds, "SELECT * FROM feedStore")

	fmt.Printf("Error: %+v\n", err3)
	fmt.Printf("Feeds: %+v\n", feeds)

}
