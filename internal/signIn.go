package internal

import (
	"errors"
	"log"
	"net/http"
)

func (is *InternalService) signInHandler(data interface{}, meta interface{}) (interface{}, int, error) {
	metaMap, ok := meta.(map[string]interface{})
	if !ok {
		log.Println("Invalid data format in signInHandler")

		return NewErrorResponse(
			"InvalidDataFormatError",
			"DFE_01",
			"Invalid data format",
		), http.StatusBadRequest, errors.New("invalid data format")
	}

	ip, ok := metaMap["ip"].(string)
	if !ok {
		log.Println("Invalid data format in signInHandler")

		return NewErrorResponse(
			"InvalidDataFormatError",
			"DFE_01",
			"Invalid data format",
		), http.StatusBadRequest, errors.New("invalid data format")
	}

	if is.isFlooder(ip) {
		log.Println("Flood protection in signInHandler")

		return NewErrorResponse(
			"FloodError",
			"DFE_07",
			"Flood protection",
		), http.StatusBadRequest, nil
	}

	dataMap, ok := data.(map[string]interface{})
	if !ok {
		log.Println("Invalid data format in signInHandler")
		is.FloodAdd(ip)

		return NewErrorResponse(
			"InvalidDataFormatError",
			"DFE_01",
			"Invalid data format",
		), http.StatusBadRequest, errors.New("invalid data format")
	}

	// Validate required fields and formats
	rules := map[string]interface{}{
		"login":    "required",
		"password": "required",
	}
	errs := is.Validate.ValidateMap(dataMap, rules)
	if len(errs) > 0 {
		log.Println("Validation error in signInHandler:", errs)
		is.FloodAdd(ip)

		return createErrorResponse(errs), http.StatusBadRequest, errors.New("not valid data")
	}

	// Hash and salt the password
	hashedPassword := is.hashAndSaltPassword(dataMap["password"].(string))

	// Check if user exists and password matches
	user, err := is.UsersRepository.GetUserByLoginAndPassword(dataMap["login"].(string), hashedPassword)
	if err != nil {
		is.FloodAdd(ip)
		return NewErrorResponse(
			"UserNotFoundError",
			"UNF_01",
			"User not found or password is incorrect",
		), http.StatusBadRequest, nil
	}

	// Generate access token and refresh token
	accessTokens, err := is.generateAccessTokens(user)
	if err != nil {
		log.Println("Cannot generate tokens, err:", err)
		return NewErrorResponse(
			"ServerError",
			"SVE_06",
			"Internal server error",
		), http.StatusInternalServerError, err
	}

	refreshToken, err := is.generateRefreshToken(user)

	if err != nil {
		log.Println("Cannot generate refresh token, err:", err)
		return NewErrorResponse(
			"ServerError",
			"SVE_06",
			"Internal server error",
		), http.StatusInternalServerError, err
	}

	user.HashedPassword = "hidden"

	// Return the tokens
	return NewOkResponse(map[string]interface{}{
		"user":         user,
		"accessTokens": accessTokens,
		"refreshToken": refreshToken,
	})
}
