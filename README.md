# rss-feeder
A simple RSS Client Server application

### Dev Notes

Just some note from me playing around

* So I learned that in normal RSS readers, my site only displays the descriptiong too. Typically,
  the reader just goes off of the rss, and my site doesn't have a <content> tag of any sort. I need
  to add that. I think we should be fine with getting the infor strictly from the feed. Maybe we can
  figure out a linking system in the client parts that can go to that link if desired.

* As far as reading html, `w3m -dump` is really good about dumping the formatted text, which might be what
  we want to connect to. 
