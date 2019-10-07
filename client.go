package requests

import (
	"time"

	"github.com/valyala/fasthttp"
)

var (
	client *fasthttp.Client
)

func init() {
	client = &fasthttp.Client{
		DialDualStack:             config.Client.DialDualStack,
		MaxConnsPerHost:           config.Client.MaxConnsPerHost,
		MaxIdleConnDuration:       config.Client.MaxIdleConnDuration * time.Second,
		MaxIdemponentCallAttempts: config.Client.MaxIdemponentCallAttempts,
		ReadBufferSize:            config.Client.ReadBufferSize,
		WriteBufferSize:           config.Client.WriteBufferSize,
		ReadTimeout:               config.Client.ReadTimeout * time.Second,
		WriteTimeout:              config.Client.WriteTimeout * time.Second,
	}
}

func Request2(method, url string, args ...interface{}) (resp *Response, err error) {
	req := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(req)

	req.SetRequestURI(url)
	req.Header.SetMethod(method)

	params := []map[string]string{}
	for _, arg := range args {
		switch a := arg.(type) {
		case Headers:
			// arg is Header , set to request header
			for k, v := range a {
				req.Header.Set(k, v)
			}
		case Params:
			// arg is "GET" params
			// ?title=website&id=1860&from=login
			params = append(params, a)
		case Auth:
			// a{username,password}
			// req.httpreq.SetBasicAuth(a[0], a[1])
		}
	}

	return
}
