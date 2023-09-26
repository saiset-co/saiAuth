package storage

import (
	"encoding/json"
	"fmt"
)

type SaiStorageSaveRequest struct {
	Collection string      `json:"collection"`
	Data       interface{} `json:"data"`
}

type SaiStorageUpdateRequest struct {
	Collection string      `json:"collection"`
	Select     interface{} `json:"select"`
	Data       interface{} `json:"data"`
}

type SaiStorageUpsertRequest = SaiStorageUpdateRequest

type SaiStorageRemoveRequest struct {
	Collection string      `json:"collection"`
	Select     interface{} `json:"select"`
}

type SaiStorageChangeResponse struct {
	Status string `json:"Status"`
}

func (saiStorage *SaiStorage) Save(request interface{}) (*SaiStorageChangeResponse, error) {
	return saiStorage.makeChangeRequest("save", request)
}

func (saiStorage *SaiStorage) Update(request interface{}) (*SaiStorageChangeResponse, error) {
	return saiStorage.makeChangeRequest("update", request)
}

func (saiStorage *SaiStorage) Upsert(request interface{}) (*SaiStorageChangeResponse, error) {
	return saiStorage.makeChangeRequest("upsert", request)
}

func (saiStorage *SaiStorage) Remove(request interface{}) (*SaiStorageChangeResponse, error) {
	return saiStorage.makeChangeRequest("remove", request)
}

func (saiStorage *SaiStorage) makeChangeRequest(method string, request interface{}) (*SaiStorageChangeResponse, error) {

	requestBody, err := json.Marshal(request)

	if err != nil {
		return nil, fmt.Errorf("failed to marshal request body: %v", err)
	}

	// Make the request
	response, err := saiStorage.makeRequest(method, requestBody)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %v", err)
	}
	defer response.Body.Close()

	// Parse the response body into the struct
	var result SaiStorageChangeResponse
	err = json.NewDecoder(response.Body).Decode(&result)
	if err != nil {
		return nil, fmt.Errorf("failed to parse response body: %v", err)
	}

	// Return the parsed results
	return &result, nil
}
