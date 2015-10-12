package main

import (
	// "flag"
	"fmt"
	"github.com/DBarney/golang-crawl/output"
	"github.com/DBarney/golang-crawl/pipeline"
	"github.com/DBarney/golang-crawl/process"
	"os"
	"os/signal"
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

	steps := []pipeline.Handler{
		process.FetchUrl,
		process.ParseXML,
		process.CompileNodeInfo,
		results.AddPage,
		process.FilterLinks,
	}
	unique := pipeline.NewPipeline(source, 1, results.IsUnique)
	rest := pipeline.NewPipeline(unique.Output(), 4, steps...)

	pending := 0
	for _, arg := range os.Args[1:] {
		pending++
		source <- arg
	}

	halt := make(chan os.Signal, 0)
	finish := make(chan interface{}, 0)
	signal.Notify(halt, os.Interrupt)
	// catch errors from all points in the pipeline
	for pending > 0 {
		select {
		case <-unique.Err():
			//nothing, as we already have this one
		case err := <-rest.Err():
			panic(err)
		case <-halt:
			fmt.Println("waiting for current jobs to finish...")
			close(finish)
		case out, open := <-rest.Output():
			if !open {
				break
			}
			newLinks := out.([]string)
			// we don't want to block the pipeline so we do this in a goroutine
			go func() {
				for _, link := range newLinks {
					select {
					case <-finish:
						pending--
					case source <- link:
					}
				}
			}()
			pending += len(newLinks)
		}
		pending--
	}
	close(source)
	results.Dump("tree")
}
