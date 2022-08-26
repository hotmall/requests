package requests

import (
	"io"
	"time"
)

// Header represent http header
type Header map[string]string

// Params represent http query params
type Params map[string]string

// Data represent http post form body
type Data map[string]string

// MultiForm represent multipart form
type MultiForm map[string]io.Reader

// Auth represent http auth, {username, password}
type Auth []string

// JSON represent http post json body
type JSON string

// Option reprents http request options
type Option struct {
	AllowRedirects                    bool
	Timeout                           time.Duration
	RequestHeaderNoDefaultContentType bool
}

// Response represents HTTP response.
type Response struct {
	noCopy     noCopy //nolint:unused,structcheck
	StatusCode int    // e.g. 200
	Header     Header
	body       []byte
}

// Content of the response, in bytes.
func (r *Response) Content() []byte {
	return r.body
}

// Text of the response, in unicode
func (r *Response) Text() string {
	return string(r.body)
}
