package output

import (
	"github.com/DBarney/golang-crawl/process"
	"testing"
)

var (
	store = &storage{}
)

func TestDropFilter(t *testing.T) {
	filter := store.FilterLinks("me.com")

	drop := &process.Page{
		Links: []string{
			// different domain
			"what.com",
			"//what.com",
			"//what.com/",
			"http://what.com",
			"http://what.com/",
			"https://what.com",
			"https://what.com/",
		},
	}

	result, err := filter(drop)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
	links := result.([]string)
	for _, link := range links {
		t.Logf("link was not dropped: %v", link)
		t.Fail()
	}
}

func TestKeepFilter(t *testing.T) {
	filter := store.FilterLinks("me.com")
	keep := &process.Page{
		Links: []string{
			"me.com/index.html",
			"http://me.com/",
			"https://me.com/",
			"//me.com/",
			"/index.html",
			"/",
			"about.com",
		},
		Url: "http://me.com/current",
	}

	result, err := filter(keep)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
	links := result.([]string)

	if len(links) != len(keep.Links) {
		t.Logf("dropped a link it shouldn't have %v %v", links, keep.Links)
		t.FailNow()
	}

	// I also need to verify that the urls are expanded correctly
	expandedLinks := []string{
		"http://me.com/index.html",
		"http://me.com/",
		"https://me.com/",
		"http://me.com/",
		"http://me.com/index.html",
		"http://me.com/",
		"http://me.com/current/about.com",
	}

	for idx, test := range expandedLinks {
		if links[idx] != test {
			t.Logf("mismatched urls %v %v", links[idx], test)
			t.Fail()
		}
	}
}
