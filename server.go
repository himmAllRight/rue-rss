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

type simpleReturn struct {
	Success string
}

type requestValues struct {
	URL      string
	Category string
}

// Takes a FeedItem object and returns it as a JSON []byte
func feedItemJSON(feedItem FeedItem) []byte {
	b, err := json.Marshal(feedItem)
	checkErrJustLog(err)
	return b
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

func readApiRequest(r *http.Request) requestValues {
	decoder := json.NewDecoder(r.Body)
	var userRequest requestValues
	err := decoder.Decode(&userRequest)
	if err != nil {
		panic(err)
	}
	defer r.Body.Close()
	log.Println(userRequest)
	return userRequest
}


////////////////////////////
//// Handler Functions ////
///////////////////////////

// Adds a new feed to the DB.
func addFeedHandler(d withDB) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		decoder := json.NewDecoder(r.Body)
		var t requestValues
		err := decoder.Decode(&t)
		if err != nil {
			panic(err)
		}
		defer r.Body.Close()
		log.Println(t)
		addFeedSource(t.URL, t.Category, d.db)
		returnStruct := simpleReturn{Success:"True"}
		fmt.Printf("Return Struct: %s\n", returnStruct)
		json.NewEncoder(w).Encode(returnStruct)
		
	})
}

// Edits feed source category
func editFeedCategory(d withDB) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userArgs := readApiRequest(r)
		editFeedSourceCat(userArgs.URL, userArgs.Category, d.db)
		fmt.Fprintf(w, "Success! The feed has been added\n")
	})
}

// Gets the FeedSource items from the feedStore
func getFeedStoreDataHandler(d withDB) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		feedStore, _ := getFeedStoreData(d.db)
		// TODO: How do we want to handle a no match? Should we just return an empty reponse?
		json.NewEncoder(w).Encode(feedStore)
	})
}

// Adds a new feed to the DB.
func deleteFeedHandler(d withDB) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userArgs := readApiRequest(r)
		deleteFeedSource(userArgs.URL, d.db)
		fmt.Fprintf(w, "Success! The feed has been removed\n")
	})
}

// Updates all the feed sources in feedStore table
func updateAllFeedsHandler(d withDB) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		updateAllFeedSources(d.db)
		fmt.Fprintf(w, "Success! All feed sources have been updated.\n")
	})
}

// Updates all the feed sources in feedStore table
func getFeedItemDataHandler(d withDB) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userArgs := readApiRequest(r)
		feedItem, _ := getFeedItemData(userArgs.URL, d.db)
		// TODO: How do we want to handle a no match? Should we just return an empty reponse?
		json.NewEncoder(w).Encode(feedItem)
	})
}

// Marks the feed item as read or unread
func markFeedItemReadHandler(readValue int, d withDB) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userArgs := readApiRequest(r)
		markReadValue(userArgs.URL, readValue, d.db)
		//	json.NewEncoder(w).Encode(feedItem)
	})
}

func getAllFeedData(d withDB) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		userArgs := readApiRequest(r)
		allFeedData, _ := getAllFeedItemData(userArgs.URL, d.db)
		json.NewEncoder(w).Encode(allFeedData)
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
	h.Handle("/edit-category", withLog(editFeedCategory(withDB(d))))
	h.Handle("/get-feedstore", withLog(getFeedStoreDataHandler(withDB(d))))
	h.Handle("/update-all-feeds", withLog(updateAllFeedsHandler(withDB(d))))
	h.Handle("/get-feeditem-data", withLog(getFeedItemDataHandler(withDB(d))))
	h.Handle("/get-all-feeditem-data", withLog(getAllFeedData(withDB(d))))
	h.Handle("/mark-read", withLog(markFeedItemReadHandler(1, withDB(d))))
	h.Handle("/mark-unread", withLog(markFeedItemReadHandler(0, withDB(d))))

	h.HandleFunc("/", noMatchHandler) // No Match condition

	err := http.ListenAndServe(":8080", h)
	log.Fatal(err)

}
