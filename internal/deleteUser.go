package internal

import (
	"github.com/saiset-co/sai-storage-mongo/external/adapter"
	"log"
	"net/http"
)

func (is *InternalService) deleteUsersHandler(data interface{}, meta interface{}) (interface{}, int, error) {
	req, ok := data.(map[string]interface{})
	if !ok || len(req) < 1 {
		return NewErrorResponse(
			"InvalidDataFormatError",
			"DFE_01",
			"Invalid data format",
		), http.StatusBadRequest, nil
	}

	updateReq := adapter.Request{
		Method: "delete",
		Data: adapter.DeleteRequest{
			Collection: "users",
			Select:     req,
		},
	}

	_, err := is.Storage.Send(updateReq)
	if err != nil {
		log.Println(
			"Cannot remove user, err:", err)
		return NewErrorResponse(
			"ServerError",
			"SVE_06",
			"Internal server error",
		), http.StatusInternalServerError, nil
	}

	return NewOkResponse("Users removed successfully")
}
