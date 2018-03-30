[![Go Report Card](https://goreportcard.com/badge/github.com/himmAllRight/rue-rss)](https://goreportcard.com/report/github.com/himmAllRight/rue-rss)

# Rue RSS
A simple RSS Client Server application. The Rue server component will runs in the background on the user's computer, or on a remote server locations. Client applications can connect to the server using a simple API. This creates a simple solution for maintaining rss feeds in a central location, and viewing them in all sorts of reader clients (web, mobile, command line).

### Dev Notes

Server: The main server component is written in go. Currently in Development.

#### TODO

- [X] Write DB component
- [X] Feed Scraper Component
    - [x] Basic API Setup
    - [ ] Implement API Calls for client applications
        - [X] Add Feed Source
        - [X] Delete Feed Source
        - [ ] Edit Feed Source [Category]
        - [X] Update All Feed Sources
        - [X] Get Feed Item Data
        - [ ] Mark Feed Item Read
        - [ ] Mark Feed Item Unread
        ...
    - [ ] Standardize all API returns
    
    ... And more I can't think of at the top of my head