package process

import (
	"net/http"
	"net/url"
)

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

func FilterLinks(baseUrl string) func(interface{}) (interface{}, error) {
	return func(job interface{}) (interface{}, error) {
		return []string{}, nil
	}
}
