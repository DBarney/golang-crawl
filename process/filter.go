package process

import (
	"net/url"
	"path/filepath"
	"regexp"
)

var (

	// these don't describe completely what a valid link looks like,
	// but used in order they can narrow down what neeeds to be done to a link
	complete         = regexp.MustCompile("^https?://[^/]+/")
	completeHostOnly = regexp.MustCompile("^https?://[^/]+")
	missingSchema    = regexp.MustCompile("^//[^/]+/")
	hostOnly         = regexp.MustCompile("^//[^/]+")
	absolute         = regexp.MustCompile("^/")
	relative         = regexp.MustCompile("^[^/]")
)

func FilterLinks(job interface{}) (interface{}, error) {
	page := job.(*Page)
	links := make([]string, 0)
	for _, link := range page.Links {
		if newLink, sameSite := expandUrl(link, page.Url); sameSite {
			links = append(links, newLink)
		}
	}
	return links, nil
}

func sameSite(first *url.URL, second *url.URL) bool {

	switch {
	case second == nil:
		return true
	default:
		return first.Host == second.Host
	}
}

func expandUrl(link string, url *url.URL) (_ string, same bool) {
	bytes := []byte(link)
	switch {
	case complete.Match(bytes):
		link = link
	case completeHostOnly.Match(bytes):
		link = link + "/"
	case missingSchema.Match(bytes):
		link = "http:" + link
	case hostOnly.Match(bytes):
		link = "http:" + link + "/"
	case absolute.Match(bytes):
		link = "http://" + url.Host + link
	case relative.Match(bytes):
		link = url.String() + "/" + link
	default:
		return "", false
	}
	testUrl, err := url.Parse(link)
	if err != nil {
		return "", false
	}
	return clean(testUrl), sameSite(testUrl, url)
}

func clean(url *url.URL) string {
	url.Fragment = ""
	url.Path = filepath.Clean(url.Path)
	return url.String()
}
