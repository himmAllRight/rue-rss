package main

import (
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

var debug = true

// Prints only if debug global is true
func debugPrint(str string) {
	if debug {
		println(str)
	}
}

// Check error returns
func checkErrFatal(err error) {
	if err != nil {
		log.Fatal("ERROR:", err)
	}
}

// Check error returns, just log if error
func checkErrJustLog(err error) {
	if err != nil {
		log.Printf("ERROR: %s\n", err)
	}
}

func main() {
	debugPrint("Load Config")
	loadConfig()

	debugPrint("Initializing DB")
	db := initDB()

	debugPrint("Adding Feed Sources from config")
	configFeedSources := configFeedSources()
	for _, configFeedSource := range configFeedSources {
		debugPrint(fmt.Sprintf("ConfigFeedSource: Category: %s URL: %s", configFeedSource.Category, configFeedSource.Feedurl))
		addFeedSource(configFeedSource.Feedurl, configFeedSource.Category, db)
	}

	updateAllFeedSources(db)
	debugPrint("hey its working.\n")
	startServer(db)
}
