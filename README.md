# requests

Go language implements a http client library like python requests, based on fasthttp.

## etc/conf/http.json

```json
{
    "dialDualStack": false,
    "maxConnsPerHost": 512,
    "maxIdleConnDuration": 10,
    "maxIdemponentCallAttempts": 5,
    "readBufferSize": 0,
    "writeBufferSize": 0,
    "readTimeout": 10,
    "writeTimeout": 10,
    "maxResponseBodySize": 100
}
```

```go
// Attempt to connect to both ipv4 and ipv6 addresses if set to true.
//
// This option is used only if default TCP dialer is used,
// i.e. if Dial is blank.
//
// By default client connects only to ipv4 addresses,
// since unfortunately ipv6 remains broken in many networks worldwide :)
DialDualStack bool

// Maximum number of connections per each host which may be established.
//
// DefaultMaxConnsPerHost(512) is used if not set.
MaxConnsPerHost int

// Idle keep-alive connections are closed after this duration.
//
// By default idle connections are closed
// after DefaultMaxIdleConnDuration(10 seconds).
MaxIdleConnDuration time.Duration

// Maximum number of attempts for idempotent calls
//
// DefaultMaxIdemponentCallAttempts(5) is used if not set.
MaxIdemponentCallAttempts int

// Per-connection buffer size for responses' reading.
// This also limits the maximum header size.
//
// Default buffer size(4096) is used if 0.
ReadBufferSize int

// Per-connection buffer size for requests' writing.
//
// Default buffer size(4096) is used if 0.
WriteBufferSize int

// Maximum duration for full response reading (including body).
//
// By default response read timeout is unlimited.
ReadTimeout time.Duration

// Maximum duration for full request writing (including body).
//
// By default request write timeout is unlimited.
WriteTimeout time.Duration

// Maximum response body size.
//
// The client returns ErrBodyTooLarge if this limit is greater than 0
// and response body is greater than the limit.
//
// By default response body size is unlimited.
MaxResponseBodySize int
```

## Post multipart form example

```go
package main

import (
   "fmt"
   "os"
   "strings"

   "github.com/hotmall/requests"
   "gitlab.com/hotmall/portmatrix"
)

func mustOpen(f string) *os.File {
    r, err := os.Open(f)
    if err != nil {
        panic(err)
    }
    return r
}

func main() {
    fmt.Println("hello world")
    remoteURL := portmatrix.MEDIA_HOST + "/media/v1/avatars"

    r := mustOpen("default.png")
    defer r.Close()

    mf := requests.MultiForm{
       "file": r,
       "id":   strings.NewReader("1111111111111111111"),
    }

    resp, err := requests.Upload(remoteURL, mf)
    if err != nil {
        panic(err)
    }
    fmt.Println("resp status is", resp.StatusCode)
    fmt.Println("resp is ", resp.Text())
}
```
