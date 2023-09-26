package internal

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/Limpid-LLC/go-auth/internal/entities"
)

const placeholder = "$"

func (is InternalService) generateAccessTokens(user *entities.User) ([]entities.AccessToken, error) {
	roles := user.Roles
	var tokenPermissions []entities.TokenPermission
	roles = append(roles, is.DefaultRole)

	for _, role := range roles {
		token, err := generateRandomToken(32)
		if err != nil {
			return nil, err
		}

		// Generate token permissions for each role
		for _, permission := range role.Permissions {
			// Generate exp time (adding one day to the current time)
			expiredAt := time.Now().Add(is.TokenExpirations.AccessToken).Unix()
			requiredParams, err := is.replacePlaceholders(permission.RequiredParams, user)
			if err != nil {
				return nil, err
			}
			restrictedParams, err := is.replacePlaceholders(permission.RestrictedParams, user)
			if err != nil {
				return nil, err
			}

			tokenPermission := entities.TokenPermission{
				Token:                      token,
				UserID:                     user.InternalId,
				Type:                       role.Type,
				StoID:                      role.StoID,
				ExpiredAt:                  expiredAt,
				RoleInternalID:             role.InternalID,
				PermissionMicroservice:     permission.Microservice,
				PermissionMethod:           permission.Method,
				PermissionRequiredParams:   requiredParams,
				PermissionRestrictedParams: restrictedParams,
			}

			tokenPermissions = append(tokenPermissions, tokenPermission)
		}
	}

	err := is.TokenPermissionsRepository.SaveTokenPermissions(tokenPermissions)

	if err != nil {
		return nil, err
	}

	var accessTokens []entities.AccessToken

	for _, tokenPermission := range tokenPermissions {
		tokenPermissionData := tokenPermission.CreateAccessToken()

		if !is.tokenPermissionExists(accessTokens, tokenPermissionData) {
			accessTokens = append(accessTokens, tokenPermissionData)
		}
	}

	return accessTokens, nil
}
func (is InternalService) tokenPermissionExists(tokens []entities.AccessToken, tokenToSearch entities.AccessToken) bool {
	for _, tokenPermission := range tokens {
		if tokenPermission.Token == tokenToSearch.Token &&
			tokenPermission.RoleId == tokenToSearch.RoleId &&
			tokenPermission.StoID == tokenToSearch.StoID &&
			tokenPermission.Type == tokenToSearch.Type {

			return true
		}
	}

	return false
}

func (is InternalService) replacePlaceholders(params []entities.Params, user *entities.User) ([]entities.Params, error) {
	for i, param := range params {
		for j, value := range param.Values {
			if len(value) > 2 && value[:1] == placeholder {
				replace, err := getEntityValue(user, value[1:])
				if err != nil {
					return nil, err
				}
				r, ok := replace.(string)
				if !ok {
					return nil, errors.New("value should be string")
				}
				params[i].Values[j] = r
			}
		}
	}

	return params, nil
}

func (is InternalService) checkHandler(data interface{}, meta interface{}) (interface{}, int, error) {

	encodedData, err := json.Marshal(data)

	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	var request Request

	err = json.Unmarshal(encodedData, &request)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	res, err := is.check(request)

	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	if !res {
		return NewErrorResponse(
			"PermissionDeniedError",
			"PDE_01",
			"Permission denied for: "+request.Microservice+"method: "+request.Method,
		), http.StatusForbidden, nil
	}

	return NewOkResponse("Ok")
}

func (is InternalService) check(request Request) (bool, error) {
	data := request.Data.(map[string]interface{})
	token, ok := data["token"].(string)
	if !ok {
		return false, errors.New("missed token")
	}

	if token == is.MasterToken {
		return true, nil
	}

	tokens, err := is.TokenPermissionsRepository.FindTokenPermissions(
		token,
		request.Microservice,
		request.Method,
	)

	if err != nil || len(tokens) == 0 {
		return false, err
	}

	return Validate(data, tokens), nil
}

func Validate(data map[string]interface{}, tokens []entities.TokenPermission) bool {
	// Check if the token has permission to access the microservice and method
	// At least one permission must be valid
	for _, permission := range tokens {
		// Validate required parameters
		if !validateRequiredParams(data, permission.PermissionRequiredParams) {
			continue
		}

		// Validate restricted parameters
		if !validateRestrictedParams(data, permission.PermissionRestrictedParams) {
			continue
		}

		// If we have validated all the required and restricted params and found no issue, return true
		return true
	}
	return false
}

func validateRequiredParams(payload map[string]interface{}, requiredParams []entities.Params) bool {
	for _, reqParam := range requiredParams {
		pathParts := strings.Split(reqParam.Param, ".")
		payloadParamValue, err := getNestedParam(payload, pathParts)
		if err != nil || payloadParamValue == nil {
			return false
		}

		// If the required param is set to All, then any value is accepted
		if reqParam.All {
			continue
		}

		// Check if the payload value matches one of the allowed values
		isMatch := false
		for _, value := range reqParam.Values {
			if value == fmt.Sprintf("%v", payloadParamValue) {
				isMatch = true
				break
			}
		}
		if !isMatch {
			return false
		}
	}
	return true
}

func validateRestrictedParams(payload map[string]interface{}, restrictedParams []entities.Params) bool {
	for _, resParam := range restrictedParams {
		pathParts := strings.Split(resParam.Param, ".")
		payloadParamValue, err := getNestedParam(payload, pathParts)
		if err != nil {
			// If restricted parameter not found in payload, it's fine
			continue
		}

		if payloadParamValue != nil {
			// If the restricted param is set to All, then any value is not accepted
			if resParam.All {
				return false
			}

			// Check if the payload value matches one of the restricted values
			for _, value := range resParam.Values {
				if value == fmt.Sprintf("%v", payloadParamValue) {
					return false
				}
			}
		}
	}
	return true
}

func getNestedParam(data map[string]interface{}, pathParts []string) (interface{}, error) {
	if len(pathParts) == 0 {
		return nil, fmt.Errorf("empty path")
	}

	value, ok := data[pathParts[0]]
	if !ok {
		return nil, fmt.Errorf("key '%s' not found", pathParts[0])
	}

	if len(pathParts) == 1 {
		// If this is the last part of the path, return the value
		return value, nil
	} else {
		// If there are more parts in the path, continue traversing
		nextData, ok := value.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("value at '%s' is not a map", pathParts[0])
		}
		return getNestedParam(nextData, pathParts[1:])
	}
}

func getEntityValue(user *entities.User, path string) (interface{}, error) {
	bytes, err := json.Marshal(user)
	if err != nil {
		return nil, err
	}

	var request map[string]interface{}

	err = json.Unmarshal(bytes, &request)
	if err != nil {
		return nil, err
	}

	value, err := getNestedParam(request, strings.Split(path, "."))
	if err != nil {
		return nil, err
	}

	return value, nil
}
