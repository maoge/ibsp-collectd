package utils

import (
	"testing"
    "fmt"
)

func TestHttpGet(t *testing.T) {
    body, err := HttpGet("http://192.168.14.206:19991/configsvr/loadServiceTopoByInstID", "INST_ID=85dc2b49-df5d-a802-e4d8-873d203015d1")
    if err == nil {
        fmt.Println(string(body))
    } else {
        t.Log(err)
    }
}

func TestHttpPost(t *testing.T) {
    body, err := HttpPost("http://192.168.14.206:19991/configsvr/loadServiceTopoByInstID", "INST_ID=85dc2b49-df5d-a802-e4d8-873d203015d1")
    if err == nil {
        fmt.Println(string(body))
    } else {
        t.Log(err)
    }
}
