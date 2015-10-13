package main

import (
	// "flag"
	"fmt"
	"github.com/DBarney/golang-crawl/output"
	"github.com/DBarney/golang-crawl/pipeline"
	"github.com/DBarney/golang-crawl/process"
	"os"
	"os/signal"
	"regexp"
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
		process.ParseHTML,
		process.CompileNodeInfo,
		results.AddPage,
		process.FilterLinks,
	}

	unique := pipeline.NewPipeline(source, 1, results.IsUnique)
	rest := pipeline.NewPipeline(unique.Output(), 4, steps...)

	pending := 0
	// this needs to be more robust.
	isDns := regexp.MustCompile("[\\w][\\d.\\w].+[\\w]")
	for _, arg := range os.Args[1:] {

		if isDns.Match([]byte(arg)) {
			pending++
			source <- "http://" + arg + "/"
			// we only want one
			break
		}
	}

	halt := make(chan os.Signal, 0)
	signal.Notify(halt, os.Interrupt)
	finish := make(chan interface{}, 0)
	go func() {
		<-halt
		fmt.Println("waiting for current jobs to finish...")
		close(finish)
	}()

	for pending > 0 {
		select {
		case <-unique.Err():
			// if we already have visited the link, we don't care about the error
		case err := <-rest.Err():
			// other errors cause the program to panic, these could be closed connections etc.
			panic(err)
		case out, open := <-rest.Output():
			if !open {
				break
			}
			newLinks := out.([]string)
			pending += len(newLinks)
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
		}
		pending--
	}
	close(source)
	results.Dump("tree")
}
