package main

import (
	// "flag"
	"fmt"
	"github.com/DBarney/golang-crawl/output"
	"github.com/DBarney/golang-crawl/pipeline"
	"github.com/DBarney/golang-crawl/process"
	"net/url"
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
	fmt.Fprintln(os.Stderr, "starting up crawler")

	results := output.NewStorage()

	steps := []pipeline.Handler{
		process.FetchUrl,
		process.ParseHTML,
		process.CompileNodeInfo,
		process.FilterLinks,
	}

	unique := pipeline.NewPipeline(source, 1, results.IsUnique, process.MaxDepth(5))
	rest := pipeline.NewPipeline(unique.Output(), 4, steps...)

	pending := 0
	for _, arg := range os.Args[1:] {
		URL, err := url.Parse(arg)
		if err == nil && (URL.Scheme == "http" || URL.Scheme == "https") {
			pending++
			if URL.Path == "" {
				URL.Path = "/"
			}
			source <- URL.String()
			break
		}
	}
	if pending == 0 {
		fmt.Fprintln(os.Stderr, "a valid http url was not provided")
		return
	}

	halt := make(chan os.Signal, 0)
	signal.Notify(halt, os.Interrupt)
	finish := make(chan interface{}, 0)
	go func() {
		<-halt
		fmt.Fprintln(os.Stderr, "waiting for current jobs to finish...")
		close(finish)
	}()

	for pending > 0 {
		select {
		case <-unique.Err():
			// if we already have visited the link, we don't care about the error
		case err := <-rest.Err():
			// other errors cause the program to exit, these could be closed connections etc.
			fmt.Fprintln(os.Stderr, "unable to continue: ", err)
			return
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
	err := results.Dump("dot", os.Stdout)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Unable to write:", err)
	}
}
