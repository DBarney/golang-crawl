package output

import (
	"errors"
	"fmt"
	"github.com/DBarney/golang-crawl/process"
	"net/url"
)

type (
	storage struct {
		store map[string]*process.Page
	}
)

var (
	Exists = errors.New("url has been visitred already")
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
	baseUrl := job.(string)
	_, exists := store.store[baseUrl]
	if exists {
		return nil, Exists
	}
	fmt.Println("fetching", baseUrl)
	URL, err := url.Parse(baseUrl)
	if err != nil {
		return nil, err
	}
	page := &process.Page{
		Url: URL,
	}
	store.store[baseUrl] = page
	return page, nil
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
