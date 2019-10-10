package requests

import "net/http"

// Header represent http header
type Header map[string]string

// Params represent http query params
type Params map[string]string

// Data represent http post form body
type Data map[string]string

// Files represent post files, [name]filename
type Files map[string]string

// Auth represent http auth, {username, password}
type Auth []string

// JSON represent http post json body
type JSON string

// Request1 represents HTTP request.
type Request1 struct {
	noCopy  noCopy
	Debug   bool
	httpreq *http.Request
	Header  *http.Header
	Cookies []*http.Cookie
}

// Response1 represents HTTP response.
type Response1 struct {
	noCopy     noCopy
	StatusCode int // e.g. 200
	Header     Header
	body       []byte
}
