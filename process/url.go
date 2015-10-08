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
	}
	return page, nil
}
