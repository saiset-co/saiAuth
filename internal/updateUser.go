package internal

import (
	"github.com/Limpid-LLC/go-auth/internal/storage"
	"github.com/Limpid-LLC/go-auth/logger"
	"go.uber.org/zap"
	"net/http"
)

func (is *InternalService) updateUserHandler(data interface{}, meta interface{}) (interface{}, int, error) {
	req, ok := data.(map[string]interface{})
	if !ok {
		return NewErrorResponse(
			"InvalidDataFormatError",
			"DFE_01",
			"Invalid data format",
		), http.StatusBadRequest, nil
	}

	selectData, ok := req["Select"].(map[string]interface{})
	if !ok || len(selectData) <= 0 {
		return NewErrorResponse(
			"MissingDataError",
			"MDE_03",
			"Missing Select data",
		), http.StatusBadRequest, nil
	}
	updateData, ok := req["Data"].(map[string]interface{})
	if !ok || len(updateData) <= 0 {
		return NewErrorResponse(
			"MissingDataError",
			"MDE_03",
			"Missing User data",
		), http.StatusBadRequest, nil
	}

	if is.userDataExists(updateData, selectData) {
		return NewErrorResponse(
			"InvalidDataFormatError",
			"RFE_04",
			"User with new data already exists",
		), http.StatusBadRequest, nil
	}

	// Restrict "___password" field
	if _, ok := updateData["___password"]; ok {
		return NewErrorResponse(
			"RestrictedFieldError",
			"RFE_02",
			"Restricted field",
		), http.StatusBadRequest, nil
	}

	// If password is provided, hash it
	if password, ok := updateData["password"]; ok {
		updateData["___password"] = is.hashAndSaltPassword(password.(string))
		delete(updateData, "password")
	}

	//todo: move to repository
	updateReq := storage.SaiStorageUpdateRequest{
		Collection: "users",
		Select:     selectData,
		Data:       updateData,
	}

	_, err := is.Storage.Update(updateReq)
	if err != nil {
		logger.Logger.Error("Cannot update user data", zap.Error(err))
		return NewErrorResponse(
			"ServerError",
			"SVE_06",
			"Internal server error",
		), http.StatusInternalServerError, nil
	}

	return NewOkResponse("User data updated successfully")
}

func (is *InternalService) userDataExists(userData, selectData map[string]interface{}) bool {
	req := storage.SaiStorageGetRequest{
		Collection: "users",
	}
	logger.Logger.Debug("userDataExists", zap.Any("userData", userData))

	filter := make(map[string]interface{})
	for key, value := range selectData {
		filter[key] = map[string]interface{}{
			"$ne": value,
		}
	}

	if email, ok := userData["email"].(string); ok {
		filterEmail := copyMap(filter)
		filterEmail["email"] = email
		req.Select = filterEmail
		res, err := is.Storage.Get(req)
		logger.Logger.Debug("userDataExists", zap.Any("res", res))
		if err != nil {
			logger.Logger.Error("Can't get existing user from the DB", zap.Error(err))
			return true
		}
		if len(res.Result) > 0 {
			return true
		}
	}

	if phone, ok := userData["phone"].(string); ok {
		filterPhone := copyMap(filter)
		filterPhone["phone"] = phone
		req.Select = filterPhone

		res, err := is.Storage.Get(req)
		logger.Logger.Debug("userDataExists", zap.Any("res", res))
		if err != nil {
			logger.Logger.Error("Can't get existing user from the DB", zap.Error(err))
			return true
		}
		if len(res.Result) > 0 {
			return true
		}
	}

	return false
}

func copyMap(input map[string]interface{}) map[string]interface{} {
	copy := make(map[string]interface{})
	for key, value := range input {
		if innerMap, ok := value.(map[string]interface{}); ok {
			copy[key] = copyMap(innerMap)
		} else {
			copy[key] = value
		}
	}
	return copy
}
