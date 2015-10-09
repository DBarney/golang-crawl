# golang-crawl ![Build Status Image](https://travis-ci.org/DBarney/golang-crawl.svg)

This mainly was written as a toy to visualize a website.

Glang-Scrape is a small web crawlr that uses goroutines to fetch webpages in parallel. The results of the crawl can be displayed in a few different formats.

tree like format
```
golang-crawl test.com -pattern=test.com -depth=10 -results=tree
test.com
|-- index.html
|-- css
|   |-- test.css
|   `-- fonts
|       |-- TimesNewRoman.font
|-- js
|   `-- all.js
```

interactive format:
```
golang-crawl test.com -pattern=test.com -depth=10 -results=interactive
please wait...
to see results visit http://127.0.0.1:4000/
type ^c to exit
```

There are a few different config options available that can be used to change how golang-crawl works.

| Parameter Name | Description | Default |
| --- | --- | --- |
| pattern | Only crawl from hosts matching the specified regex pattern | The url passed as the first parameter|
| depth | Only follow links that are within this number of links of the first base page | 10 |
| results | Change how results are displayed | tree |


## notes
- I need to have a site map generated with links to other pages and static resources that are loaded.

## TODO
- be able to scan css for links to static resources as well
- actually match the parameters in this README.md
