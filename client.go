package requests

import (
	"bytes"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/valyala/fasthttp"
)

const (
	// VERSION represent hotmall/requests version
	VERSION = "0.1"
)

var (
	client *fasthttp.Client
)

func init() {
	client = &fasthttp.Client{
		Name:                      "Hotmall Go Requests " + VERSION,
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

// Request sends a http request
func Request(method, url string, args ...interface{}) (resp *Response, err error) {
	req := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(req)

	opts := buildRequest(req, method, url, args...)

	response := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(response)

	err = doRequestTimeout(req, response, opts.AllowRedirects, opts.Timeout)
	if err != nil {
		return
	}

	resp, err = buildResponse(response)
	return
}

func buildRequest(req *fasthttp.Request, method, url string, args ...interface{}) Option {
	var opts = Option{
		AllowRedirects: false,
		Timeout:        0,
	}
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
		case Option:
			opts = arg.(Option)
		default:

		}
	}

	req.SetRequestURI(url)
	req.Header.SetMethod(method)

	return opts
}

func buildResponse(r *fasthttp.Response) (resp *Response, err error) {
	resp = &Response{
		Header: make(Header),
	}
	resp.StatusCode = r.StatusCode()

	r.Header.VisitAll(func(key, value []byte) {
		resp.Header[string(key)] = string(value)
	})

	var b []byte
	if v := r.Header.Peek(fasthttp.HeaderContentEncoding); v != nil {
		if bytes.Compare(v, []byte("gzip")) == 0 {
			b, err = r.BodyGunzip()
			if err != nil {
				return
			}
		} else if bytes.Compare(v, []byte("deflate")) == 0 {
			b, err = r.BodyInflate()
			if err != nil {
				return
			}
		} else {
			err = fmt.Errorf("Not support Content-Encoding:%s", v)
			return
		}
	} else {
		b = r.Body()
	}

	resp.body = append(resp.body, b...)
	return
}

var errorPool sync.Pool

func doRequestTimeout(req *fasthttp.Request, resp *fasthttp.Response, allowRedirects bool, timeout time.Duration) (err error) {
	if timeout <= 0 {
		return doRequestFollowRedirects(req, resp, allowRedirects)
	}

	var ch chan error
	chv := errorPool.Get()
	if chv == nil {
		chv = make(chan error, 1)
	}
	ch = chv.(chan error)

	go func() {
		err := doRequestFollowRedirects(req, resp, allowRedirects)
		ch <- err
	}()

	tc := fasthttp.AcquireTimer(timeout)
	select {
	case err = <-ch:
		errorPool.Put(chv)
	case <-tc.C:
		err = fasthttp.ErrTimeout
	}
	fasthttp.ReleaseTimer(tc)

	return
}

const maxRedirectsCount = 16

var (
	errMissingLocation  = errors.New("missing Location header for http redirect")
	errTooManyRedirects = errors.New("too many redirects detected when doing the request")
)

func doRequestFollowRedirects(req *fasthttp.Request, resp *fasthttp.Response, allowRedirects bool) (err error) {
	url := req.RequestURI()
	redirectsCount := 0
	for {
		req.SetRequestURIBytes(url)

		if err = client.Do(req, resp); err != nil {
			break
		}

		statusCode := resp.Header.StatusCode()
		if !fasthttp.StatusCodeIsRedirect(statusCode) {
			break
		}

		if !allowRedirects {
			break
		}

		redirectsCount++
		if redirectsCount > maxRedirectsCount {
			err = errTooManyRedirects
			break
		}

		location := resp.Header.Peek(fasthttp.HeaderLocation)
		if len(location) == 0 {
			err = errMissingLocation
			break
		}
		url = getRedirectURL(location)
	}

	return
}

func getRedirectURL(location []byte) []byte {
	u := fasthttp.AcquireURI()
	defer fasthttp.ReleaseURI(u)
	u.UpdateBytes(location)
	redirectURL := u.FullURI()
	return redirectURL
}
