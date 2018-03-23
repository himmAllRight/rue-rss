package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"html"
	"log"
	"net/http"
)

type testStruct struct {
	Test string
}

type feedEntry struct {
	FeedURL  string
	Category string
}

//addFeedSource
func startServer(db *sql.DB) {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, %q", html.EscapeString(r.URL.Path))
	})

	http.HandleFunc("/test", func(rw http.ResponseWriter, req *http.Request) {
		decoder := json.NewDecoder(req.Body)
		var t testStruct
		err := decoder.Decode(&t)
		if err != nil {
			panic(err)
		}
		defer req.Body.Close()
		log.Println(t.Test)
	})
	http.HandleFunc("/add-feed", func(rw http.ResponseWriter, req *http.Request) {
		decoder := json.NewDecoder(req.Body)
		var t feedEntry
		err := decoder.Decode(&t)
		if err != nil {
			panic(err)
		}
		defer req.Body.Close()
		log.Println(t)
		addFeedSource(t.FeedURL, t.Category, db)
	})

	log.Fatal(http.ListenAndServe(":8080", nil))

}
