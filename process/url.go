package process

import (
	"errors"
	"github.com/DBarney/golang-crawl/pipeline"
	"net/http"
)

var (
	DepthExceeded = errors.New("max depth exceeded")
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

func MaxDepth(depth int) pipeline.Handler {
	return func(job interface{}) (interface{}, error) {
		// page := job.(*Page)
		// if page.Depth > depth {
		// 	return nil, DepthExceeded
		// }
		return job, nil
	}
}
