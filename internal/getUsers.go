package internal

import (
	"log"
	"net/http"

	"github.com/saiset-co/sai-storage-mongo/external/adapter"
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

	getReq := adapter.Request{
		Method: "read",
		Data: adapter.ReadRequest{
			Collection: "users",
			Select:     selectData,
		},
	}

	res, err := is.Storage.Send(getReq)
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
		delete(user, "___password")
	}

	return NewOkResponse(res.Result)
}
