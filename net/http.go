package net

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
)

func HttpGet(url string, optHeader ...map[string]string) ([]byte, error) {

	if url == "" {
		return nil, errors.New("url is empty")
	}

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	return ioutil.ReadAll(resp.Body)
}

func HttpPostJson(url string, postBody interface{}) ([]byte, error) {

	if url == "" {
		return nil, errors.New("url is empty")
	}

	var jsonStr []byte
	if postBody == nil {
		jsonStr = make([]byte, 0)
	} else {
		switch v := postBody.(type) {
		case string:
			jsonStr = []byte(v)
		default:
			var e error
			jsonStr, e = json.Marshal(postBody)
			if e != nil {
				return nil, e
			}
		}
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonStr))
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	return ioutil.ReadAll(resp.Body)
}
