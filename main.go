package main

import (
	// "flag"
	"fmt"
	"github.com/DBarney/golang-crawl/output"
	"github.com/DBarney/golang-crawl/pipeline"
	"github.com/DBarney/golang-crawl/process"
	"os"
)

func init() {
	// flag.StringVar(&pattern, "pattern", "", "only urls that match the pattern will be crawled")
	// flag.StringVar(&output, "output", "tree", "the results of the crawl are displayed in this format")
	// flag.IntVar(&depth, "depth", 10, "the depth that the crawler is allowed to dispaly")
}

// usage golang-crawl url [pattern] [depth]
func main() {
	source := make(chan interface{}, 0)

	// flag.Parse()
	fmt.Println("starting up crawler")

	results := output.NewStorage()

	unique := pipeline.NewPipeline(source, results.IsUnique)
	urlFetcher := pipeline.NewPipeline(unique.Output(), process.FetchUrl)
	xmlParser := pipeline.NewPipeline(urlFetcher.Output(), process.ParseXML)
	documentCompiler := pipeline.NewPipeline(xmlParser.Output(), process.CompileNodeInfo)
	storage := pipeline.NewPipeline(documentCompiler.Output(), results.AddPage)
	links := pipeline.NewPipeline(storage.Output(), results.FilterLinks())

	pending := 0
	for _, arg := range os.Args[1:] {
		pending++
		source <- arg
	}

	// catch errors from all points in the pipeline
	for pending > 0 {
		select {
		case <-unique.Err():
			//nothing, as we already have this one
		case err := <-urlFetcher.Err():
			panic(err)
		case err := <-xmlParser.Err():
			panic(err)
		case err := <-documentCompiler.Err():
			panic(err)
		case err := <-storage.Err():
			panic(err)
		case err := <-links.Err():
			panic(err)
		case out := <-links.Output():
			newLinks := out.([]string)
			// we don't want to block the pipeline so we do this in a goroutine
			go func() {
				for _, link := range newLinks {
					source <- link
				}
			}()
			pending += len(newLinks)
		}
		pending--
	}
	results.Dump("tree")
}
