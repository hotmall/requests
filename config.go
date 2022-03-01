package requests

import (
	"time"
)

var (
	config httpConfig
)

type httpConfig struct {
	Client clientConfig `json:"client"`
}

type clientConfig struct {
	DialDualStack             bool          `json:"dialDualStack"`
	MaxConnsPerHost           int           `josn:"maxConnsPerHost"`
	MaxIdleConnDuration       time.Duration `json:"maxIdleConnDuration"`
	MaxIdemponentCallAttempts int           `json:"maxIdemponentCallAttempts"`
	ReadBufferSize            int           `json:"readBufferSize"`
	WriteBufferSize           int           `json:"writeBufferSize"`
	ReadTimeout               time.Duration `json:"readTimeout"`
	WriteTimeout              time.Duration `json:"writeTimeout"`
	MaxResponseBodySize       int           `json:"maxResponseBodySize"`
}

func (c *clientConfig) Initial() {
	c.DialDualStack = false
	c.MaxConnsPerHost = 512
	c.MaxIdleConnDuration = 10 * time.Second
	c.MaxIdemponentCallAttempts = 5
	c.ReadBufferSize = 4096
	c.WriteBufferSize = 4096
	c.ReadTimeout = 10 * time.Second
	c.WriteTimeout = 10 * time.Second
	c.MaxResponseBodySize = 2 * 1024 * 1024
}

// func loadConfig(filename string, c *httpConfig) (err error) {
// 	if !fileExist(filename) {
// 		err = fmt.Errorf("file not found, file=%s", filename)
// 		return
// 	}
// 	bytes, err := ioutil.ReadFile(filename)
// 	if err != nil {
// 		return
// 	}
// 	if err = json.Unmarshal(bytes, c); err != nil {
// 		return
// 	}
// 	return
// }

// func fileExist(file string) bool {
// 	_, err := os.Stat(file)
// 	if err != nil && os.IsNotExist(err) {
// 		return false
// 	}
// 	return true
// }

func init() {
	config.Client.Initial()
	// prefix := commandline.PrefixPath()
	// filename := prefix + "/etc/conf/http.json"
	// err := loadConfig(filename, &config)
	// if err != nil {
	// 	fmt.Printf("Load config file(%s) fail, err = %s\n", filename, err.Error())
	// }
}
