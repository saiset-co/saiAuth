package internal

import (
	"errors"
	"log"
	"net/http"
)

func (is *InternalService) checkOTPCodeHandler(data interface{}, meta interface{}) (interface{}, int, error) {
	dataMap, ok := data.(map[string]interface{})
	if !ok {
		log.Println("Invalid data format in checkOTPCodeHandler")
		return NewErrorResponse(
			"InvalidDataFormatError",
			"DFE_01",
			"Invalid data format",
		), http.StatusBadRequest, errors.New("invalid data format")
	}

	// Validate required fields and formats
	rules := map[string]interface{}{
		"otp_code": "required",
	}

	errs := is.Validate.ValidateMap(dataMap, rules)
	if len(errs) > 0 {
		log.Println("Validation error in checkOTPCodeHandler:", errs)
		return createErrorResponse(errs), http.StatusBadRequest, errors.New("not valid data")
	}

	phone, phonePresent := dataMap["phone"]
	email, emailPresent := dataMap["email"]

	if (phonePresent && isFieldEmpty(phone)) || (emailPresent && isFieldEmpty(email)) {
		return NewErrorResponse(
			"InvalidDataFormatError",
			"DFE_01",
			"Invalid data format",
		), http.StatusBadRequest, errors.New("phone or email is required")
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

	return NewOkResponse("OTP code is valid")
}

func isFieldEmpty(value interface{}) bool {
	switch v := value.(type) {
	case string:
		return v == ""
	default:
		return true
	}
}
