package process

import (
	"net/url"
	"testing"
)

func TestFilterDrop(t *testing.T) {
	url, err := url.Parse("http://me.com/current")
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
	drop := &Page{
		Links: []string{
			// different domain
			"//what.com/",
			"//what.com",
			"http://what.com/",
			"http://what.com",
			"https://what.com/",
		},
		Url: url,
	}

	result, err := FilterLinks(drop)
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

func TestFilterKeep(t *testing.T) {
	url, err := url.Parse("http://me.com/current")
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
	keep := &Page{
		Links: []string{
			"http://me.com/",
			"http://me.com//",
			"http://me.com/././../../../../../../index.html",
			"https://me.com/",
			"//me.com/",
			"//me.com",
			"/index.html",
			"/",
			"about.com",
			"about.com#",
		},
		Url: url,
	}

	result, err := FilterLinks(keep)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
	links := result.([]string)

	if len(links) != len(keep.Links) {
		t.Logf("dropped a link it shouldn't have %v %v", links, keep.Links)
		t.FailNow()
	}

	expandedLinks := []string{
		"http://me.com/",
		"http://me.com/",
		"http://me.com/index.html",
		"https://me.com/",
		"http://me.com/",
		"http://me.com/",
		"http://me.com/index.html",
		"http://me.com/",
		"http://me.com/current/about.com",
		"http://me.com/current/about.com",
	}

	for idx, test := range expandedLinks {
		if links[idx] != test {
			t.Logf("link was expanded incorrectly %v %v", links[idx], test)
			t.Fail()
		}
	}
}
