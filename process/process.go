package process

import (
	"golang.org/x/net/html"
	"net/http"
	"net/url"
)

type (
	Page struct {
		Links  []string
		Assets []string
		Root   *html.Node
		Res    *http.Response
		Url    *url.URL
	}
)
