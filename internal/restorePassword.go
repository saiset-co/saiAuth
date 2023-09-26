package internal

import (
	"encoding/json"
	"fmt"
	"github.com/Limpid-LLC/go-auth/internal/storage"
	"log"
	"net/http"
	"time"
)

type RestorePasswordRequest struct {
	Phone    string `json:"phone"`
	Email    string `json:"email" validate:"required_without=Phone"`
	OtpCode  string `json:"otp_code"`
	Password string `json:"password"`
}

func (is *InternalService) restorePasswordHandler(data interface{}, meta interface{}) (interface{}, int, error) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		log.Println("Invalid data format in signUpHandler")
		return NewErrorResponse(
			"InvalidDataFormatError",
			"DFE_01",
			"Invalid data format",
		), http.StatusInternalServerError, fmt.Errorf("error Marshal" + err.Error())
	}

	var request RestorePasswordRequest
	err = json.Unmarshal(jsonData, &request)
	if err != nil {
		return NewErrorResponse(
			"InvalidDataFormatError",
			"DFE_01",
			"Invalid data format",
		), http.StatusInternalServerError, fmt.Errorf("error Unmarshal" + err.Error())
	}

	errs := is.Validate.Struct(request)
	if errs != nil {
		log.Printf("Validation errors: %v", errs)

		return NewErrorResponse(
			"InvalidDataFormatError",
			"VLE_03",
			"Validation error.",
		), http.StatusBadRequest, errs
	}

	// Check OTP code
	ok := is.checkOTPCodeRestorePassword(request)
	if !ok {
		return NewErrorResponse(
			"OTPError",
			"OPE_05",
			"Invalid OTP code",
		), http.StatusBadRequest, nil
	}

	user, err := is.UsersRepository.GetUserByPhoneOrEmail(request.Phone, request.Email)
	if err != nil {
		return NewErrorResponse(
			"UserExistsError",
			"UEE_04",
			"User with this phone or email not found",
		), http.StatusBadRequest, err
	}

	user.HashedPassword = is.hashAndSaltPassword(request.Password)

	err = is.UsersRepository.UpdateUser(user)
	if err != nil {
		log.Println("Cannot update user data, err:", err)
		return NewErrorResponse(
			"ServerError",
			"SVE_06",
			"Internal server error",
		), http.StatusInternalServerError, nil
	}

	// Remove OTP code from storage
	_, err = is.Storage.Remove(storage.SaiStorageRemoveRequest{
		Collection: "otpCodes",
		Select: map[string]interface{}{
			"code":  request.OtpCode,
			"phone": request.Phone,
			"email": request.Email,
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

	return NewOkResponse("Restore password successfully")
}

func (is *InternalService) checkOTPCodeRestorePassword(request RestorePasswordRequest) bool {
	otpRes, err := is.Storage.Get(storage.SaiStorageGetRequest{
		Collection: "otpCodes",
		Select: map[string]interface{}{
			"code":  request.OtpCode,
			"phone": request.Phone,
			"email": request.Email,
			"expired_at": map[string]interface{}{
				"$gte": time.Now(),
			},
		},
	})
	return err == nil && len(otpRes.Result) > 0
}
