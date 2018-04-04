package utils

import (
    "fmt"
    "strings"
    "io/ioutil"
    "net/http"
)

func HttpGet(url, params string) ([]byte, error) {
    client := &http.Client{}
    req, err := http.NewRequest("GET", url+"?"+params, nil)
    req.Header.Set("Content-Type", "application/text")
    // req.Header.Set("Cookie", "name=anny")

    resp, err := client.Do(req)
    defer resp.Body.Close()
    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        fmt.Println(err)
    }

    return body, err
}

func HttpPost(url, params string) ([]byte, error) {
    client := &http.Client{}
    req, err := http.NewRequest("POST", url, strings.NewReader(params))
    req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
    // req.Header.Set("Cookie", "name=anny")

    resp, err := client.Do(req)
    defer resp.Body.Close()
    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        fmt.Println(err)
    }

    return body, err
}

