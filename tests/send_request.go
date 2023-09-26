package tests

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
)

const (
	baseUrl = "http://localhost:9001"
)

type Payload struct {
	Method string                 `json:"method"`
	Data   map[string]interface{} `json:"data"`
}

func sendRequest(method string, url string, payload interface{}) (interface{}, int, error) {
	client := &http.Client{}
	jsonPayload, _ := json.Marshal(payload)
	req, err := http.NewRequest(method, url, bytes.NewBuffer(jsonPayload))
	if err != nil {
		return nil, 0, err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		return nil, 0, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, resp.StatusCode, err
	}
	var result interface{}
	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, resp.StatusCode, err
	}
	return result, resp.StatusCode, nil
}
