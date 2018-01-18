package main

import (
	"fmt"
	"github.com/k3a/html2text"
	//"io"
	"io/ioutil"
)

// Checks for errors
func check(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {
	// Reads contents of file
	dat, err := ioutil.ReadFile("/tmp/input.html")
	check(err)

	filehtml := string(dat)
	html := `<html><head><title>Test Page</title></head><body>Some <b>Content</b> <strong>here</strong></body>`

	filePlain := html2text.HTML2Text(filehtml)

	plain := html2text.HTML2Text(html)

	fmt.Println(plain)
	fmt.Println(filePlain)
}
