package main

import (
	"encoding/json"
	"fmt"
	"html"
	"log"
	"net/http"
)

type test_struct struct {
	Test string
}

type feed_entry struct {
	feedURL  string
	category string
}

//addFeedSource
func startServer() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, %q", html.EscapeString(r.URL.Path))
	})

	http.HandleFunc("/test", func(rw http.ResponseWriter, req *http.Request) {
		decoder := json.NewDecoder(req.Body)
		var t test_struct
		err := decoder.Decode(&t)
		if err != nil {
			panic(err)
		}
		defer req.Body.Close()
		log.Println(t.Test)
	})
	http.HandleFunc("/add-feed", func(rw http.ResponseWriter, req *http.Request) {
		decoder := json.NewDecoder(req.Body)
		var t feed_entry
		err := decoder.Decode(&t)
		if err != nil {
			panic(err)
		}
		defer req.Body.Close()
		log.Println(t)
	})

	log.Fatal(http.ListenAndServe(":8080", nil))

}
