package process

import (
	"net/http"
)

func FetchUrl(job interface{}) (interface{}, error) {
	url := job.(string)
	res, err := http.Get(url)
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
