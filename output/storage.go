package output

import (
	"fmt"
	"github.com/DBarney/golang-crawl/process"
	"net/url"
	"regexp"
)

type (
	storage struct {
		store map[string]*process.Page
	}
)

var (
	// these don't describe completely what a valid link looks like,
	// but used in order they can narrow down what a link is
	complete       = regexp.MustCompile("^https?://[^/]+/")
	missing_schema = regexp.MustCompile("^//[^/]+/")
	absolute       = regexp.MustCompile("^/")
	relative       = regexp.MustCompile("^[^/]")
)

func NewStorage() *storage {
	store := &storage{
		store: make(map[string]*process.Page, 1),
	}
	return store
}
func (store *storage) AddPage(job interface{}) (interface{}, error) {
	page := job.(*process.Page)
	key := page.Url.String()
	_, visited := store.store[key]
	if !visited {
		store.store[key] = page
	}

	return page, nil
}

func (store *storage) FilterLinks(pattern string) func(interface{}) (interface{}, error) {
	return func(job interface{}) (interface{}, error) {
		page := job.(*process.Page)
		links := make([]string, 0)
		for _, link := range page.Links {
			if newLink, sameSite := matchUrl(link, page.Url); sameSite {
				links = append(links, newLink)
			}
		}
		return links, nil
	}
}

func sameSite(link string, url *url.URL) bool {
	testUrl, err := url.Parse(link)
	switch {
	case err != nil:
		return false
	default:
		return testUrl.Host == url.Host
	}
}

func matchUrl(link string, url *url.URL) (string, bool) {
	// url.Parse doesn't really work here, so I need to regex my way to victory
	bytes := []byte(link)
	switch {
	case complete.Match(bytes):
		return link, sameSite(link, url)
	case missing_schema.Match(bytes):
		full := "http:" + link
		return full, sameSite(full, url)
	case absolute.Match(bytes):
		return "http://" + url.Host + link, true
	case relative.Match(bytes):
		return url.String() + "/" + link, true
	default:
		return "", false
	}
}

func (store *storage) Dump(method string) {
	for url, page := range store.store {
		fmt.Printf("%v\n", url)
		for _, resource := range page.Assets {
			fmt.Printf("  | %v\n", resource)
		}
		for _, resource := range page.Links {
			fmt.Printf("  * %v\n", resource)
		}
	}
}
