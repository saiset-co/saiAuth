package internal

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/Limpid-LLC/go-auth/internal/entities"
	"github.com/saiset-co/sai-storage-mongo/external/adapter"

	"github.com/go-playground/validator/v10"
)

type Request struct {
	Microservice string      `json:"microservice"`
	Method       string      `json:"method"`
	Data         interface{} `json:"data"`
}

func (is *InternalService) createRoleHandler(data interface{}, meta interface{}) (interface{}, int, error) {

	jsonData, err := json.Marshal(data)
	if err != nil {
		panic(err)
	}

	var role entities.Role
	err = json.Unmarshal(jsonData, &role)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	err = validator.New().Struct(role)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	req := adapter.Request{
		Method: "create",
		Data: adapter.CreateRequest{
			Collection: "roles",
			Documents:  []interface{}{role},
		},
	}

	_, err = is.Storage.Send(req)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	return NewOkResponse("Role created successfully")
}

func (is *InternalService) updateRolesHandler(data interface{}, meta interface{}) (interface{}, int, error) {
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
			"MDE_04",
			"Missing Data for update",
		), http.StatusBadRequest, nil
	}

	updateReq := adapter.Request{
		Method: "update",
		Data: adapter.UpdateRequest{
			Collection: "roles",
			Select:     selectData,
			Document:   map[string]interface{}{"$set": updateData},
		},
	}

	getReq := adapter.Request{
		Method: "read",
		Data: adapter.ReadRequest{
			Collection: "roles",
			Select:     selectData,
		},
	}

	_, err := is.Storage.Send(updateReq)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	rolesData, err := is.Storage.Send(getReq)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	var roles []entities.Role
	jsonData, err := json.Marshal(rolesData.Result)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}
	err = json.Unmarshal(jsonData, &roles)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	for _, role := range roles {
		// remove related tokens

		err = is.TokenPermissionsRepository.RemoveTokenPermissionsByRoleInternalID(role.InternalID)

		if err != nil {
			return nil, http.StatusInternalServerError, err
		}

		// update related users
		err = is.updateRoleInUsers(&role)
		if err != nil {
			return nil, http.StatusInternalServerError, err
		}
	}

	return NewOkResponse("Role updated successfully")
}

func (is *InternalService) deleteRolesHandler(data interface{}, meta interface{}) (interface{}, int, error) {
	dataMap, ok := data.(map[string]interface{})
	if !ok || len(dataMap) <= 2 {
		log.Println("Invalid data format in deleteRoleHandler")
		return NewErrorResponse(
			"InvalidDataFormatError",
			"DFE_03",
			"Invalid data format",
		), http.StatusBadRequest, nil
	}

	getReq := adapter.Request{
		Method: "read",
		Data: adapter.ReadRequest{
			Collection: "roles",
			Select:     dataMap,
		},
	}

	req := adapter.Request{
		Method: "delete",
		Data: adapter.DeleteRequest{
			Collection: "roles",
			Select:     dataMap,
		},
	}

	rolesData, err := is.Storage.Send(getReq)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	if len(rolesData.Result) == 0 {
		return nil, http.StatusInternalServerError, errors.New("no roles to delete by the request")
	}

	var roles []entities.Role
	jsonData, err := json.Marshal(rolesData.Result)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}
	err = json.Unmarshal(jsonData, &roles)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	_, err = is.Storage.Send(req)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	for _, role := range roles {
		err = is.TokenPermissionsRepository.RemoveTokenPermissionsByRoleInternalID(role.InternalID)
		if err != nil {
			return nil, http.StatusInternalServerError, err
		}

		// remove role from users
		err = is.deleteRoleFromUsers(role.InternalID)
		if err != nil {
			return nil, http.StatusInternalServerError, err
		}
	}

	return NewOkResponse("Role deleted successfully")
}

func (is *InternalService) attachRole(userID string, roleID string) error {
	// Fetch the user
	user, err := is.UsersRepository.GetUserByID(userID)
	if err != nil {
		return err
	}

	getReq := adapter.Request{
		Method: "read",
		Data: adapter.ReadRequest{
			Collection: "roles",
			Select: map[string]interface{}{
				"internal_id": roleID,
			},
		},
	}

	// Fetch the role
	rolesData, err := is.Storage.Send(getReq)
	if err != nil {
		return err
	}

	if len(rolesData.Result) == 0 {
		return errors.New("role not found")
	}

	var roles []entities.Role
	jsonData, err := json.Marshal(rolesData.Result)
	if err != nil {
		return err
	}
	err = json.Unmarshal(jsonData, &roles)
	if err != nil {
		return err
	}

	// Attach the role to the user
	user.AddRole(roles[0])

	// Update the user
	err = is.UsersRepository.UpdateUser(user)
	if err != nil {
		return err
	}

	return nil
}

func (is *InternalService) detachRole(userID string, roleID string) error {
	// Fetch the user
	user, err := is.UsersRepository.GetUserByID(userID)
	if err != nil {
		return err
	}

	// Detach the role from the user
	user.DeleteRole(roleID)

	// Update the user
	err = is.UsersRepository.UpdateUser(user)
	if err != nil {
		return err
	}

	return nil
}

func (is *InternalService) updateRoleInUsers(role *entities.Role) error {
	// Fetch users who have this role
	users, err := is.UsersRepository.GetUsersByRole(role.InternalID)
	if err != nil {
		return err
	}

	// For each user, update the role
	for _, user := range users {
		user.UpdateRole(*role)

		err = is.UsersRepository.UpdateUser(&user)
		if err != nil {
			return err
		}
	}

	return nil
}

func (is *InternalService) deleteRoleFromUsers(internalID string) error {
	// Fetch users who have this role
	users, err := is.UsersRepository.GetUsersByRole(internalID)
	if err != nil {
		return err
	}

	// For each user, remove the role
	for _, user := range users {
		user.DeleteRole(internalID)

		err = is.UsersRepository.UpdateUser(&user)
		if err != nil {
			return err
		}
	}

	return nil
}

func (is *InternalService) attachRoleHandler(data interface{}, meta interface{}) (interface{}, int, error) {
	dataMap, ok := data.(map[string]interface{})
	if !ok {
		log.Println("Invalid data format in attachRoleHandler")
		return NewErrorResponse(
			"InvalidDataFormatError",
			"DFE_04",
			"Invalid data format",
		), http.StatusBadRequest, nil
	}

	userID, ok := dataMap["user_id"].(string)
	if !ok {
		return NewErrorResponse(
			"InvalidUserIDError",
			"IUE_01",
			"Invalid user ID",
		), http.StatusBadRequest, nil
	}

	roleID, ok := dataMap["role_id"].(string)
	if !ok {
		return NewErrorResponse(
			"InvalidRoleIDError",
			"IRE_01",
			"Invalid role ID",
		), http.StatusBadRequest, nil
	}

	err := is.attachRole(userID, roleID)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	return NewOkResponse("Role attached successfully")
}

func (is *InternalService) detachRoleHandler(data interface{}, meta interface{}) (interface{}, int, error) {
	dataMap, ok := data.(map[string]interface{})
	if !ok {
		log.Println("Invalid data format in detachRoleHandler")
		return NewErrorResponse(
			"InvalidDataFormatError",
			"DFE_05",
			"Invalid data format",
		), http.StatusBadRequest, nil
	}

	userID, ok := dataMap["user_id"].(string)
	if !ok {
		return NewErrorResponse(
			"InvalidUserIDError",
			"IUE_02",
			"Invalid user ID",
		), http.StatusBadRequest, nil
	}

	roleID, ok := dataMap["role_id"].(string)
	if !ok {
		return NewErrorResponse(
			"InvalidRoleIDError",
			"IRE_02",
			"Invalid role ID",
		), http.StatusBadRequest, nil
	}

	err := is.detachRole(userID, roleID)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	return NewOkResponse("Role detached successfully")
}

func (is *InternalService) getRolesHandler(data interface{}, meta interface{}) (interface{}, int, error) {
	selectData, ok := data.(map[string]interface{})
	delete(selectData, "service_station_id")
	if !ok {
		log.Println("Invalid data format in getRolesHandler")
		return NewErrorResponse(
			"InvalidDataFormatError",
			"DFE_01",
			"Invalid data format",
		), http.StatusBadRequest, nil
	}

	req := adapter.Request{
		Method: "read",
		Data: adapter.ReadRequest{
			Collection: "roles",
			Select:     selectData,
		},
	}

	res, err := is.Storage.Send(req)
	if err != nil {
		log.Println("Cannot Get from sai storage, err:", err)
		return NewErrorResponse(
			"ServerError",
			"SVE_06",
			"Internal server error",
		), http.StatusInternalServerError, err
	}

	return NewOkResponse(res.Result)
}
