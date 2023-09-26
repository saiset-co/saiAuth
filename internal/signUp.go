package internal

import (
	"errors"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/Limpid-LLC/go-auth/internal/storage"
)

//todo: rewrite this for using struct User

func (is *InternalService) signUpHandler(data interface{}, meta interface{}) (interface{}, int, error) {
	dataMap, ok := data.(map[string]interface{})
	if !ok {
		log.Println("Invalid data format in signUpHandler")
		return NewErrorResponse(
			"InvalidDataFormatError",
			"DFE_01",
			"Invalid data format",
		), http.StatusBadRequest, errors.New("invalid data format")
	}

	// Check for restricted fields
	ok, errResp := is.checkRestrictedFields(dataMap)
	if !ok {
		return errResp, http.StatusBadRequest, errors.New("restricted fields")
	}

	// Validate required fields and formats
	rules := map[string]interface{}{
		"email":    "required|email",
		"phone":    "required",
		"otp_code": "required",
		"password": "required",
	}
	errs := is.Validate.ValidateMap(dataMap, rules)
	if len(errs) > 0 {
		log.Println("Validation error in signUpHandler:", errs)
		return createErrorResponse(errs), http.StatusBadRequest, errors.New("not valid data")
	}

	// Check OTP code
	ok = is.checkOTPCode(dataMap)
	if !ok {
		return NewErrorResponse(
			"OTPError",
			"OPE_05",
			"Invalid OTP code",
		), http.StatusBadRequest, nil
	}

	// Check if user exists
	ok = is.checkUserExistence(dataMap)
	if !ok {
		return NewErrorResponse(
			"UserExistsError",
			"UEE_04",
			"User with this email or phone already exists",
		), http.StatusBadRequest, nil
	}

	userData, _ := dataMap["data"].(interface{})

	user := is.createUser(dataMap["email"].(string), dataMap["phone"].(string), dataMap["password"].(string), userData)

	err := is.UsersRepository.CreateUser(user)

	if err != nil {
		log.Println("Cannot Save to sai storage, err:", err)
		return NewErrorResponse(
			"ServerError",
			"SVE_06",
			"Internal server error",
		), http.StatusInternalServerError, err
	}

	// Remove OTP code from storage
	_, err = is.Storage.Remove(storage.SaiStorageRemoveRequest{
		Collection: "otpCodes",
		Select: map[string]interface{}{
			"code":  dataMap["otp_code"],
			"phone": dataMap["phone"],
		},
	})

	if err != nil {
		log.Println("Cannot remove OTP code from storage, err:", err)
		return NewErrorResponse(
			"ServerError",
			"SVE_06",
			"Internal server error",
		), http.StatusInternalServerError, err
	}

	return NewOkResponse("User registered successfully")
}

func (is *InternalService) checkRestrictedFields(dataMap map[string]interface{}) (bool, *ErrorResponse) {
	for key := range dataMap {
		if strings.HasPrefix(key, "___") {
			log.Println("Attempt to use restricted field in signUpHandler:", key)
			errResp := NewErrorResponse(
				"RestrictedFieldError",
				"RFE_02",
				"Restricted field",
			)

			return false, &errResp
		}
	}
	return true, nil
}

func (is *InternalService) checkOTPCode(dataMap map[string]interface{}) bool {
	otpRes, err := is.Storage.Get(storage.SaiStorageGetRequest{
		Collection: "otpCodes",
		Select: map[string]interface{}{
			"code": dataMap["otp_code"],
			"$or": []map[string]interface{}{
				{"email": dataMap["email"]},
				{"phone": dataMap["phone"]},
			},
			"expired_at": map[string]interface{}{
				"$gte": time.Now(),
			},
		},
	})
	return err == nil && len(otpRes.Result) > 0
}

func (is *InternalService) checkUserExistence(dataMap map[string]interface{}) bool {
	res, err := is.Storage.Get(storage.SaiStorageGetRequest{
		Collection: "users",
		Select: map[string]interface{}{
			"$or": []map[string]interface{}{
				{"email": dataMap["email"]},
				{"phone": dataMap["phone"]},
			},
		},
	})
	return err != nil || len(res.Result) == 0
}
