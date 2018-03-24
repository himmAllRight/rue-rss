package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"html"
	"log"
	"net/http"
)

type appHandler func(http.ResponseWriter, *http.Request) (int, error)

// Request Structs
type testStruct struct {
	Test string
}

type feedEntry struct {
	FeedURL  string
	Category string
}

// Handler Functions

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

// How can I pass the DB object to the handler function?
func addFeedHandler(rw http.ResponseWriter, req *http.Request) {
	decoder := json.NewDecoder(req.Body)
	var t feedEntry
	err := decoder.Decode(&t)
	if err != nil {
		panic(err)
	}
	defer req.Body.Close()
	log.Println(t)
	//addFeedSource(t.FeedURL, t.Category, db)
	fmt.Fprintf(rw, "Success! The feed has been added\n")
}

func withLog(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("Before")
		next.ServeHTTP(w, r)
		log.Println("After")
	})
}

// Server
func startServer(db *sql.DB) {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, %q", html.EscapeString(r.URL.Path))
	})
	http.HandleFunc("/feed-store", func(w http.ResponseWriter, r *http.Request) {
		feedStore := getStoreFeeds(db)
		fmt.Fprintf(w, "%q", feedStore)
	})
	http.HandleFunc("/test", apiHandler)
	http.HandleFunc("/add-feed", func(w http.ResponseWriter, r *http.Request) {
		decoder := json.NewDecoder(r.Body)
		var t feedEntry
		err := decoder.Decode(&t)
		if err != nil {
			panic(err)
		}
		defer r.Body.Close()
		log.Println(t)
		addFeedSource(t.FeedURL, t.Category, db)
		fmt.Fprintf(w, "Success! The feed has been added\n")
	})

	log.Fatal(http.ListenAndServe(":8080", nil))

}
