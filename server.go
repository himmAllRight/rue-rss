package main

import (
	"encoding/json"
	"fmt"
	"html"
	"log"
	"net/http"

	"github.com/jmoiron/sqlx"
)

// Request Structs
type testStruct struct {
	Test string
}

type feedEntry struct {
	FeedURL  string
	Category string
}

///////////////////////////////////
//// Handler Wrapper Functions ////
///////////////////////////////////

// Struct and method to pass db into handlers
type withDB struct {
	db *sqlx.DB
}

func (d withDB) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	println("db handler")
}

// Wrapper that logs before and after handler. Might be used later.
func withLog(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("Before")
		h.ServeHTTP(w, r) // Call orig
		log.Println("After")
	})
}

////////////////////////////
//// Handler Functions ////
///////////////////////////

// Generalized handler functions
func apiHandler(rq http.ResponseWriter, req *http.Request) {
	decoder := json.NewDecoder(req.Body)
	var t testStruct
	err := decoder.Decode(&t)
	if err != nil {
		panic(err)
	}
	defer req.Body.Close()
	log.Println(t.Test)
}

// Adds a new feed to the DB.
func addFeedHandler(d withDB) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		decoder := json.NewDecoder(r.Body)
		var t feedEntry
		err := decoder.Decode(&t)
		if err != nil {
			panic(err)
		}
		defer r.Body.Close()
		log.Println(t)
		addFeedSource(t.FeedURL, t.Category, d.db)
		fmt.Fprintf(w, "Success! The feed has been added\n")
	})
}

// Adds a new feed to the DB.
func deleteFeedHandler(d withDB) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		decoder := json.NewDecoder(r.Body)
		var t feedEntry
		err := decoder.Decode(&t)
		if err != nil {
			panic(err)
		}
		defer r.Body.Close()
		log.Println(t)
		deleteFeedSource(t.FeedURL, d.db)
		fmt.Fprintf(w, "Success! The feed has been removed\n")
	})
}

// Updates all the feed sources in feedStore table
func updateAllFeedsHandler(d withDB) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		decoder := json.NewDecoder(r.Body)
		var t feedEntry
		err := decoder.Decode(&t)
		if err != nil {
			panic(err)
		}
		defer r.Body.Close()
		log.Println(t)
		updateAllFeedSources(d.db)
		fmt.Fprintf(w, "Success! All feed sources have been updated.\n")
	})
}

func noMatchHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Undefined request url, %q\n", html.EscapeString(r.URL.Path))
}

////////////////////
//// Server Run ////
////////////////////

func startServer(db *sqlx.DB) {
	h := http.NewServeMux()
	// Create a db handler obj to pass the db pointer around
	d := withDB{db}

	// Handler conditions
	h.Handle("/add-feed", withLog(addFeedHandler(withDB(d))))
	h.Handle("/delete-feed", withLog(deleteFeedHandler(withDB(d))))
	h.Handle("/update-all-feeds", withLog(updateAllFeedsHandler(withDB(d))))

	h.HandleFunc("/test", apiHandler) // Simple API test
	h.HandleFunc("/", noMatchHandler) // No Match condition

	err := http.ListenAndServe(":8080", h)
	log.Fatal(err)

}
