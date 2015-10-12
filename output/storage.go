package output

import (
	"errors"
	"fmt"
	"github.com/DBarney/golang-crawl/process"
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
