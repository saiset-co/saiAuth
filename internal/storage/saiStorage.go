package storage

import (
	"bytes"
	"fmt"
	"net/http"
)

type SaiStorage struct {
	Url   string
	Token string
}

func (s *SaiStorage) makeRequest(method string, requestBody []byte) (*http.Response, error) {
	// Define the request URLx
	url := s.Url + "/" + method

	//println(requestBody)
	println(string(requestBody))

	// Create a new POST request with the request body
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	// Add the Token header to the request
	req.Header.Set("Token", s.Token)

	// Send the request and get the response
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %v", err)
	}

	return resp, nil
}
