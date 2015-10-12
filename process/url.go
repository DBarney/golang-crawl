package process

import (
	"net/http"
	"net/url"
	"regexp"
	"strings"
)

type (
	patTrans struct {
		pattern   *regexp.Regexp
		transform func(string, *url.URL) string
	}
)

var (

	// these don't describe completely what a valid link looks like,
	// but used in order they can narrow down what neeeds to be done to a link
	check = []patTrans{
		patTrans{regexp.MustCompile("^https?://[^/]+/"), complete},
		patTrans{regexp.MustCompile("^https?://[^/]+"), completeHostOnly},
		patTrans{regexp.MustCompile("^//[^/]+/"), missingSchema},
		patTrans{regexp.MustCompile("^//[^/]+"), hostOnly},
		patTrans{regexp.MustCompile("^/"), absolute},
		patTrans{regexp.MustCompile("^[^/]"), relative},
	}
)

func complete(match string, url *url.URL) string         { return match }
func completeHostOnly(match string, url *url.URL) string { return match + "/" }
func missingSchema(match string, url *url.URL) string    { return "http:" + match }
func hostOnly(match string, url *url.URL) string         { return "http:" + match + "/" }
func absolute(match string, url *url.URL) string         { return "http://" + url.Host + match }
func relative(match string, url *url.URL) string         { return url.String() + "/" + match }

func FetchUrl(job interface{}) (interface{}, error) {
	stringUrl := job.(string)
	url, err := url.Parse(stringUrl)
	if err != nil {
		return nil, err
	}
	res, err := http.Get(stringUrl)
	if err != nil {
		return nil, err
	}
	page := &Page{
		Res: res,
		Url: url,
	}
	return page, nil
}

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

func sameSite(link string, url *url.URL) bool {
	testUrl, err := url.Parse(link)
	switch {
	case err != nil:
		return false
	default:
		return testUrl.Host == url.Host
	}
}

func expandUrl(link string, url *url.URL) (string, bool) {
	bytes := []byte(link)
	for _, pair := range check {
		if pair.pattern.Match(bytes) {
			link = pair.transform(link, url)
			return strip(link), sameSite(link, url)
		}
	}
	return "", false
}

func strip(url string) string {
	noFragment := strings.SplitN(url, "#", 2)[0]
	return noFragment
}
