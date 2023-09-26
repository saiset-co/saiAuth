package internal

import (
	"github.com/Limpid-LLC/go-auth/internal/storage"
	"log"
	"net/http"
)

func (is *InternalService) getUsersHandler(data interface{}, meta interface{}) (interface{}, int, error) {
	selectData, ok := data.(map[string]interface{})
	if !ok {
		log.Println("Invalid data format in getUsersHandler")
		return NewErrorResponse(
			"InvalidDataFormatError",
			"DFE_01",
			"Invalid data format",
		), http.StatusBadRequest, nil
	}

	// Fetch users from SaiStorage
	getReq := storage.SaiStorageGetRequest{
		Collection: "users",
		Select:     selectData,
	}
	res, err := is.Storage.Get(getReq)
	if err != nil {
		log.Println("Cannot Get from sai storage, err:", err)
		return NewErrorResponse(
			"ServerError",
			"SVE_06",
			"Internal server error",
		), http.StatusInternalServerError, err
	}

	// Remove hashed passwords from the response
	for _, user := range res.Result {
		userMap, ok := user.(map[string]interface{})
		if ok {
			delete(userMap, "___password")
		}
	}

	return NewOkResponse(res.Result)
}
