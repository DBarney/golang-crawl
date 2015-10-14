package output

import (
	"errors"
	"fmt"
	"github.com/DBarney/golang-crawl/process"
	"io"
	"net/url"
	"os"
)

type (
	storage struct {
		store map[string]*process.Page
		order int
	}
	exporter func(map[string]*process.Page, io.Writer) error

	reader struct{}
)

var (
	Exists          = errors.New("url has been visitred already")
	UnknownExporter = errors.New("unknown export type")
	exporters       = map[string]exporter{
		"dot": dot,
	}
)

func NewStorage() *storage {
	store := &storage{
		store: make(map[string]*process.Page, 1),
		order: 0,
	}
	return store
}

func (store *storage) IsUnique(job interface{}) (interface{}, error) {
	baseUrl := job.(string)
	_, exists := store.store[baseUrl]
	if exists {
		return nil, Exists
	}
	fmt.Fprintln(os.Stderr, "fetching", baseUrl)
	URL, err := url.Parse(baseUrl)
	if err != nil {
		return nil, err
	}
	page := &process.Page{
		Url: URL,
		Id:  store.order,
	}
	store.order++
	store.store[baseUrl] = page
	return page, nil
}

func dot(data map[string]*process.Page, dst io.Writer) error {
	dst.Write([]byte("digraph {"))
	for url, page := range data {
		fmt.Fprintf(dst, " %v -> {", page.Id)
		sep := ""
		for _, resource := range page.SameDomainLinks {
			remoteId := data[resource].Id
			if remoteId == page.Id {
				continue
			}
			fmt.Fprintf(dst, "%v%v", sep, remoteId)
			sep = " "
		}
		fmt.Fprintf(dst, "} %v [label=\"\",title=\"%v\",size=%v]", page.Id, url, len(page.Assets))
	}
	_, err := dst.Write([]byte(" }"))
	return err
}

func (store *storage) Dump(method string, dst io.Writer) error {
	exporter, exists := exporters[method]
	if !exists {
		return UnknownExporter
	}
	return exporter(store.store, dst)
}
