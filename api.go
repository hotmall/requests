package requests

// Get sends a GET request
func Get(url string, args ...interface{}) (resp *Response, err error) {
	return Request(MethodGet, url, args...)
}

// Post sends a POST request
func Post(url string, args ...interface{}) (resp *Response, err error) {
	return Request(MethodPost, url, args...)
}

// Put sends a PUT request
func Put(url string, args ...interface{}) (resp *Response, err error) {
	return Request(MethodPut, url, args...)
}

// Delete send a DELETE request
func Delete(url string, args ...interface{}) (resp *Response, err error) {
	return Request(MethodDelete, url, args...)
}

// Head send a HEAD request
func Head(url string, args ...interface{}) (resp *Response, err error) {
	return Request(MethodHead, url, args...)
}

// Patch send a PATCH request
func Patch(url string, args ...interface{}) (resp *Response, err error) {
	return Request(MethodPatch, url, args...)
}

// Options sends an OPTIONS request
func Options(url string, args ...interface{}) (resp *Response, err error) {
	return Request(MethodOptions, url, args...)
}

// Upload post multipart form request
func Upload(url string, args ...interface{}) (resp *Response, err error) {
	return Request(MethodPost, url, args...)
}
