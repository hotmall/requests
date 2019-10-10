package requests

import (
	"bytes"
	"fmt"
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

func Request2(method, url string, args ...interface{}) (resp *Response1, err error) {
	req := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(req)

	params := []map[string]string{}
	for _, arg := range args {
		switch a := arg.(type) {
		case Header:
			// arg is Header , set to request header
			for k, v := range a {
				req.Header.Set(k, v)
			}
		case Params:
			// arg is "GET" params
			// ?title=website&id=1860&from=login
			params = append(params, a)
			args := fasthttp.AcquireArgs()
			defer fasthttp.ReleaseArgs(args)

			params := arg.(Params)
			for k, v := range params {
				args.Add(k, v)
			}
			s := args.String()
			if len(s) > 0 {
				url += "?" + s
			}

		case Data:
			data := arg.(Data)
			args := req.PostArgs()
			for k, v := range data {
				args.Add(k, v)
			}
			req.SetBody(args.QueryString())
			req.Header.SetContentType("application/x-www-form-urlencoded")
		case JSON:
			req.Header.SetContentType("application/json")
			req.SetBodyString(arg.(string))
		case Auth:
			// a{username,password}
			// req.httpreq.SetBasicAuth(a[0], a[1])
		}
	}

	req.SetRequestURI(url)
	req.Header.SetMethod(method)

	response := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(response)

	err = client.Do(req, response)
	if err != nil {
		return
	}

	resp.StatusCode = response.StatusCode()
	var b []byte
	if v := response.Header.Peek(fasthttp.HeaderContentEncoding); v != nil {
		if bytes.Compare(v, []byte("gzip")) == 0 {
			b, err = response.BodyGunzip()
			if err != nil {
				return
			}
		} else if bytes.Compare(v, []byte("deflate")) == 0 {
			b, err = response.BodyInflate()
			if err != nil {
				return
			}
		} else {
			err = fmt.Errorf("Not support Content-Encoding:%s", v)
			return
		}
	} else {
		b = response.Body()
	}

	resp.body = append(resp.body, b...)
	return
}
