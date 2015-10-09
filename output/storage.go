package output

import (
	"fmt"
	"github.com/DBarney/golang-crawl/process"
)

type (
	storage struct {
		store map[string]*process.Page
	}
)

func NewStorage() *storage {
	store := &storage{
		store: make(map[string]*process.Page, 1),
	}
	return store
}
func (store *storage) AddPage(job interface{}) (interface{}, error) {
	page := job.(*process.Page)
	_, visited := store.store[page.Url]
	if !visited {
		store.store[page.Url] = page
	}

	return page, nil
}

func (store *storage) FilterLinks(pattern string) func(interface{}) (interface{}, error) {
	return func(job interface{}) (interface{}, error) {
		// I need to filter out all links not matching the pattern, also I need to fix
		// up any urls that aren't complete so that they can be requested correctly (relative etc.)
		return []string{}, nil
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
