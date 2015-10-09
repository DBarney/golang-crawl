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

	urlFetcher := pipeline.NewPipeline(source, process.FetchUrl)
	xmlParser := pipeline.NewPipeline(urlFetcher.Output(), process.ParseXML)
	documentCompiler := pipeline.NewPipeline(xmlParser.Output(), process.CompileNodeInfo)
	storage := pipeline.NewPipeline(documentCompiler.Output(), results.AddPage)
	links := pipeline.NewPipeline(storage.Output(), storage.FilterLinks("what.org"))

	pending := 0
	for _, arg := range os.Args[1:] {
		pending++
		source <- arg
	}

	// catch errors from all points in the pipeline
	for pending > 0 {
		select {
		case err := <-urlFetcher.Err():
			panic(err)
		case err := <-xmlParser.Err():
			panic(err)
		case err := <-storage.Err():
			panic(err)
		case err := <-documentCompiler.Err():
			panic(err)
		case out := <-links.Output():
			newLinks := out.([]string)
			for _, link := range newLinks {
				fmt.Printf("also need to get %v\n", link)
			}
		}
		pending--
	}
	results.Dump("tree")
}
