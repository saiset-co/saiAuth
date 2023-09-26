package internal

import (
	"github.com/Limpid-LLC/go-auth/internal/storage"
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

	updateReq := storage.SaiStorageRemoveRequest{
		Collection: "users",
		Select:     req,
	}

	_, err := is.Storage.Remove(updateReq)
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
