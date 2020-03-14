package requests

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"github.com/mallbook/commandline"
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

func loadConfig(filename string, c *httpConfig) (err error) {
	if !fileExist(filename) {
		err = fmt.Errorf("File not found, file=%s", filename)
		return
	}
	bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		return
	}
	if err = json.Unmarshal(bytes, c); err != nil {
		return
	}
	return
}

func fileExist(file string) bool {
	_, err := os.Stat(file)
	if err != nil && os.IsNotExist(err) {
		return false
	}
	return true
}

func init() {
	prefix := commandline.PrefixPath()
	filename := prefix + "/etc/conf/http.json"
	err := loadConfig(filename, &config)
	if err != nil {
		fmt.Printf("Load config file(%s) fail, err = %s\n", filename, err.Error())
	}
}
