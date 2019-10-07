package requests

import "net/http"

// Headers represent http header
type Headers map[string]string

// Params represent http query params
type Params map[string]string

// Data represent http post form body
type Data map[string]string

// Files represent post files, [name]filename
type Files map[string]string

// Auth represent http auth, {username, password}
type Auth []string

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
	noCopy  noCopy
	R       *http.Response
	content []byte
	text    string
	req     *Request1
}
