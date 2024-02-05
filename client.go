package requests

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"strings"
	"sync"
	"time"

	"github.com/valyala/fasthttp"
)

const (
	FILE_KEY        = "file"
	FILENAME_KEY    = "filename"
	defaultFilename = "hotfile"
)

var (
	client *fasthttp.Client
)

func init() {
	client = &fasthttp.Client{
		Name:                      "Hotmall Go Requests",
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
	request := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(request)

	opts, err := buildRequest(request, method, url, args...)
	if err != nil {
		return
	}
	request.Header.SetNoDefaultContentType(opts.RequestHeaderNoDefaultContentType)

	response := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(response)

	err = doRequestTimeout(request, response, opts.AllowRedirects, opts.Timeout)
	if err != nil {
		return
	}

	resp, err = buildResponse(response)
	return
}

func buildRequest(req *fasthttp.Request, method, url string, args ...interface{}) (opts Option, err error) {
	for _, arg := range args {
		switch t := arg.(type) {
		case Header:
			// arg is Header , set to request header
			h := arg.(Header)
			for k, v := range h {
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
			j := arg.(JSON)
			req.SetBodyString(string(j))
		case MultiForm:
			mf := arg.(MultiForm)
			err = buildMultiForm(req, mf)
		case Auth:
			// a{username,password}
			// req.httpreq.SetBasicAuth(a[0], a[1])
		case Option:
			opts = arg.(Option)
		default:
			err = fmt.Errorf("not support argument type:(%s)", t)
		}
	}
	req.SetRequestURI(url)
	req.Header.SetMethod(method)
	return
}

func buildMultiForm(req *fasthttp.Request, mf MultiForm) (err error) {
	var b bytes.Buffer
	w := multipart.NewWriter(&b)
	// 先解析 filename
	filename := parseFilename(mf)

	for key, f := range mf {
		var part io.Writer
		if key == FILE_KEY {
			// Add a media file
			if part, err = w.CreateFormFile(key, filename); err != nil {
				return
			}
		} else {
			// Add other fields
			if part, err = w.CreateFormField(key); err != nil {
				return
			}
		}
		if _, err = io.Copy(part, f); err != nil {
			return err
		}
	}
	// Don't forget to close the multipart writer.
	// If you don't close it, your request will be missing the terminating boundary.
	w.Close()

	contentType := w.FormDataContentType()
	req.Header.SetContentType(contentType)
	req.SetBody(b.Bytes())
	return
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
		if bytes.Equal(v, []byte("gzip")) {
			b, err = r.BodyGunzip()
			if err != nil {
				return
			}
		} else if bytes.Equal(v, []byte("deflate")) {
			b, err = r.BodyInflate()
			if err != nil {
				return
			}
		} else {
			err = fmt.Errorf("not support Content-Encoding:%s", v)
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
		if !StatusCodeIsRedirect(statusCode) {
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

// StatusCodeIsRedirect returns true if the status code indicates a redirect.
func StatusCodeIsRedirect(statusCode int) bool {
	return statusCode == StatusMovedPermanently ||
		statusCode == StatusFound ||
		statusCode == StatusSeeOther ||
		statusCode == StatusTemporaryRedirect ||
		statusCode == StatusPermanentRedirect
}

func parseFilename(mf MultiForm) string {
	filename := defaultFilename
	if r, ok := mf[FILENAME_KEY]; ok {
		if rdr, ok := r.(*strings.Reader); ok {
			b := make([]byte, rdr.Size())
			if _, err := rdr.Read(b); err == nil {
				filename = string(b)
				// 不要忘记重新 reset 以下，否则再次调用 Read 方法，会返回 io.EOF 错误
				rdr.Reset(filename)
			}
		}
	}
	return filename
}
