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

func sqlxTestMain() {
	db, err := sqlx.Connect("sqlite3", "test-db2.db")
	if err != nil {
		log.Fatalln(err)
	}

	db.MustExec(schema)

	tx := db.MustBegin()
	tx.MustExec("INSERT INTO feedStore (feedurl, category) VALUES ($1, $2)", "http://ryan.himmelwright.net/post/index.xml", "Test")
	tx.Commit()

	feeds := []FeedSource{}
	err3 := db.Select(&feeds, "SELECT * FROM feedStore")

	fmt.Printf("Error: %+v\n", err3)
	fmt.Printf("People: %+v\n", feeds)

}
