package spl

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
)

// RequestGet ..
func RequestGet(path string, out interface{}) (*http.Response, error) {
	return request(http.MethodGet, path, &out, nil)
}

func request(method string, path string, out interface{}, data interface{}) (*http.Response, error) {
	var req *http.Request
	var err error

	if data != nil {
		var bytesData []byte
		bytesData, err = json.Marshal(data)
		if err != nil {
			return nil, err
		}

		req, err = http.NewRequest(method, "https://api.servers.com/v1/"+path, bytes.NewBuffer(bytesData))
		if err != nil {
			return nil, err
		}
	} else {
		req, err = http.NewRequest(method, "https://api.servers.com/v1/"+path, nil)
		if err != nil {
			return nil, err
		}
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+Conf.JWTtoken)
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != 200 {
		return nil, fmt.Errorf("response status code: " + res.Status)
	}

	go RateLimitNotify(res)

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	err = json.Unmarshal(body, &out)
	if err != nil {
		return nil, err
	}

	return res, nil
}

// StrToInt ..
func StrToInt(s string) int {
	number, err := strconv.ParseInt(s, 10, 0)
	if err != nil {
		return -1
	}
	return int(number)
}
