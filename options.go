package requests

type Options struct {
	Headers map[string]string
	Cookies map[string]string
	Params  map[string]string
	Data    string
	Json    string
}
