package process

import (
	"net/http"
)

func FetchUrl(job interface{}) (interface{}, error) {
	page := job.(*Page)
	res, err := http.Get(page.Url.String())
	if err != nil {
		return nil, err
	}
	page.Res = res
	return page, nil
}
