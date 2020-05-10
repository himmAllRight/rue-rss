[![Go Report Card](https://goreportcard.com/badge/github.com/himmAllRight/rue-rss)](https://goreportcard.com/report/github.com/himmAllRight/rue-rss)

# Rue RSS
A simple RSS Client Server application. The Rue server component will runs in the background on the user's computer, or on a remote server locations. Client applications can connect to the server using a simple API. This creates a simple solution for maintaining rss feeds in a central location, and viewing them in all sorts of reader clients (web, mobile, command line).

## Config

A config file can be provided to provide rue with specific information. The
config should be a `yaml` file named `config.yaml`, and can be located at
either `~/.config/rue/config.yaml`, `/etc/rue/config.yaml`, or the current
directory (`./config.yaml`).

Example Config:

```yaml
db:
  src: "test-db.db"   # The realitive location to store the databse file

feed_sources:         # Defines new Feed Sources to add
  me:                 # Each sub-item defines a Feed Source Category
    # Below that, each item is a feed url in that category
    - "http://ryan.himmelwright.net/post/index.xml"
  fedora:
    - "https://fedoramagazine.org/feed/"
```

## Dev Notes

Server: The main server component is written in go. Currently in Development.

### API Requests to Server


#### Get Feedstore

```bash
curl -X POST http://localhost:8080/get-all-feeditem-data
```

#### Get All Feed Item Data

```bash
curl -X POST -d "{\"URL\":\"http://ryan.himmelwright.net/post/index.xml\"}" http://localhost:8080/get-all-feeditem-data
```

#### Get All Feed Item Data

```bash
curl -X POST -d "{\"URL\":\"http://ryan.himmelwright.net/post/index.xml\", \"Category\": \"Deveopment\"}" http://localhost:8080/add-feed
```


### Server TODO

- [X] Write DB component
- [X] Feed Scraper Component
- [x] Basic API Setup
    - [ ] Implement API Calls for client applications
        - [X] Add Feed Source
        - [X] Delete Feed Source
        - [X] Get feedStore
        - [X] Edit Feed Source [Category]
        - [X] Update All Feed Sources
        - [X] Get Feed Item Data
		- [X] Get All Feed Item Data from a Feed Source
        - [X] Mark Feed Item Read
        - [X] Mark Feed Item Unread
        ...
    - [X] Standardize all API returns
    ... And more I can't think of at the top of my head

- [ ] Setup Preferences System?
    - [X] Load Preferences
    - [ ] Create Default Preferences if Missing
    - [X] Define config options and values
    - [ ] Update system to get values from config during run
        - [X] DB file location default or from config
        - [X] Feed Sources listed in config file
        - [ ] Optional Integrations loaded if if config file
    ...
- Integrations
    - [ ] Pocket
        - [ ] Authenticate with Pocket
        - [ ] Import Feed Items into pocket list
            - [ ] Only New items in DB
            - [ ] Remain in DB if failed import
- [ ] Authentication System
    ...
- [ ] Project Build
    - [ ] Properly organize application
    - [ ] Determine build steps
    - [ ] Write build instructions

