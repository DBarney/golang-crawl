package output

import (
	"errors"
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
	Exists = errors.New("url has been visitred already")

	// these don't describe completely what a valid link looks like,
	// but used in order they can narrow down what a link is
	complete         = regexp.MustCompile("^https?://[^/]+/")
	completeHostOnly = regexp.MustCompile("^https?://[^/]+")
	missingSchema    = regexp.MustCompile("^//[^/]+/")
	hostOnly         = regexp.MustCompile("^//[^/]+")
	absolute         = regexp.MustCompile("^/")
	relative         = regexp.MustCompile("^[^/]")
)

func NewStorage() *storage {
	store := &storage{
		store: make(map[string]*process.Page, 1),
	}
	return store
}
func (store *storage) AddPage(job interface{}) (interface{}, error) {
	page := job.(*process.Page)
	store.store[page.Url.String()] = page
	return page, nil
}

func (store *storage) IsUnique(job interface{}) (interface{}, error) {
	url := job.(string)
	_, exists := store.store[url]
	if exists {
		return nil, Exists
	}
	// just a place holder
	store.store[url] = &process.Page{}
	fmt.Println("fetching", url)
	return job, nil
}

func (store *storage) FilterLinks() func(interface{}) (interface{}, error) {
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
	case completeHostOnly.Match(bytes):
		full := link + "/"
		return full, sameSite(full, url)
	case missingSchema.Match(bytes):
		full := "http:" + link
		return full, sameSite(full, url)
	case hostOnly.Match(bytes):
		full := "http:" + link + "/"
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
