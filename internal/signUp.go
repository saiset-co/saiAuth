package internal

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/Limpid-LLC/go-auth/internal/storage"
)

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

	phone, _ := dataMap["phone"].(string)
	email, _ := dataMap["email"].(string)
	otpCode, _ := dataMap["otp_code"].(string)
	password, _ := dataMap["password"].(string)

	if phone == "" && email == "" {
		log.Printf("Validation errors: phone or email required")
		return nil, http.StatusBadRequest, fmt.Errorf("phone or email required")
	}

	if (is.SmsEnabled || is.EmailEnabled) && otpCode == "" {
		log.Printf("Validation errors: otp_code is required")
		return nil, http.StatusBadRequest, fmt.Errorf("otp_code is required")
	}

	if password == "" {
		log.Printf("Validation errors: password is required")
		return nil, http.StatusBadRequest, fmt.Errorf("password is required")
	}

	if is.EmailEnabled || is.SmsEnabled {
		// Check OTP code
		ok = is.checkOTPCode(dataMap)
		if !ok {
			return NewErrorResponse(
				"OTPError",
				"OPE_05",
				"Invalid OTP code",
			), http.StatusBadRequest, nil
		}
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

	user := is.createUser(email, phone, password, userData)

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
			"code": otpCode,
			"$or": []map[string]interface{}{
				{"email": email},
				{"phone": phone},
			},
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
