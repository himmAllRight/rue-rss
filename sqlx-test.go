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

type Person struct {
	FirstName string `db:"first_name"`
	LastName  string `db:"last_name"`
	Email     string
}

type feedStore struct {
	id       int    `db:"id"`
	feedurl  string `db:"feedurl"`
	category string `db:"category"`
}

type feedData struct {
	id              int
	feedname        string
	feedurl         string
	postname        string
	posturl         string
	publishdate     string
	postdescription string
	postcontent     string
}

func sqlxTestMain() {
	db, err := sqlx.Connect("sqlite3", "test-db2.db")
	if err != nil {
		log.Fatalln(err)
	}

	db.MustExec(schema)

	tx := db.MustBegin()
	//tx.MustExec("INSERT INTO feedStore (feedurl, category) VALUES ($1, $2)", "http://ryan.himmelwright.net/post/index.xml", "Test")
	tx.MustExec("INSERT INTO person (first_name, last_name, email) VALUES ($1, $2, $3)", "Ryan", "Himmelwright", "ryan@himmelwright.net")
	tx.MustExec("INSERT INTO person (first_name, last_name, email) VALUES ($1, $2, $3)", "Rebecca", "Himmelwright", "rebecca@himmelwright.net")
	tx.NamedExec("INSERT INTO person (first_name, last_name, email) VALUES (:first_name, :last_name, :email)", &Person{"Jane", "Citizen", "jane.citzen@example.com"})
	tx.Commit()

	people := []Person{}
	err2 := tx.Select(&people, "SELECT * FROM person")
	fmt.Printf("Error: %+v\n", err2)
	fmt.Printf("People: %+v\n", people)

}
