package internal

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/saiset-co/sai-storage-mongo/external/adapter"
	"log"
	"math/rand"
	"net/http"
	"time"
)

type OTPCode struct {
	Code      string    `json:"code"`
	ExpiredAt time.Time `json:"expired_at"`
	Phone     string    `json:"phone"`
	Email     string    `json:"email"`
}

type SMSData struct {
	Phone     string                 `json:"phone"`
	Message   string                 `json:"message"`
	Template  string                 `json:"template"`
	Variables map[string]interface{} `json:"variables"`
}

type SMSRequest struct {
	Method string  `json:"method"`
	Data   SMSData `json:"data"`
}

type EmailRequest struct {
	Method string    `json:"method"`
	Data   EmailData `json:"data"`
}

type EmailData struct {
	Recipient string `json:"recipient"`
	Sender    string `json:"sender"`
	Subject   string `json:"subject"`
	Body      string `json:"body"`
}

func (is *InternalService) sendVerifyCodeHandler(data interface{}, meta interface{}) (interface{}, int, error) {
	dataMap, ok := data.(map[string]interface{})
	if !ok {
		log.Printf("Invalid data format: %v", data)
		return nil, http.StatusBadRequest, fmt.Errorf("invalid data format")
	}

	rules := map[string]interface{}{
		"phone": "required",
	}

	errs := is.Validate.ValidateMap(dataMap, rules)

	if len(errs) > 0 {
		log.Printf("Validation errors: %v", errs)
		return createErrorResponse(errs), http.StatusBadRequest, fmt.Errorf("invalid data format")
	}

	phone, _ := dataMap["phone"].(string)
	template, templateOk := dataMap["template"].(string)
	if !templateOk {
		template = ""
	}
	variables, variablesOk := dataMap["variables"].(map[string]interface{})
	if !variablesOk {
		variables = map[string]interface{}{}
	}

	// Generate a 4-digit OTP code
	code := fmt.Sprintf("%04d", rand.Intn(10000))

	fakeMode := dataMap["fake"] == is.MasterKey

	// If the "fake" key exists and matches the master key, don't send SMS and set the code to "1111"
	if fakeMode {
		code = "1111"
	}

	// Set the expiration date (e.g., 5 minutes from now)
	expDate := time.Now().Add(5 * time.Minute)

	saveReq := adapter.Request{
		Method: "create",
		Data: adapter.CreateRequest{
			Collection: "otpCodes",
			Documents: []interface{}{
				OTPCode{
					Code:      code,
					ExpiredAt: expDate,
					Phone:     phone,
				},
			},
		},
	}

	_, err := is.Storage.Send(saveReq)
	if err != nil {
		log.Printf("Error saving OTP code to SaiStorage: %v", err)
		return nil, http.StatusInternalServerError, err
	}

	log.Printf("OTP code saved to SaiStorage successfully")

	// If not a fake request, send the OTP code via SMS
	if !fakeMode {
		err = is.sendSMS(phone, "Code: "+code, template, variables)
		if err != nil {
			log.Printf("Error sending SMS: %v", err)
			return nil, http.StatusInternalServerError, err
		}
	}

	log.Printf("OTP code sent successfully")

	return NewOkResponse("OTP code sent")
}

type ResetPasswordVerifyCodeRequest struct {
	Phone     string                 `json:"phone"`
	Email     string                 `json:"email" validate:"required_without=Phone"`
	Template  string                 `json:"template"`
	Variables map[string]interface{} `json:"variables"`
	Fake      string                 `json:"fake"`
}

func (is *InternalService) sendResetPasswordVerifyCodeHandler(data interface{}, meta interface{}) (interface{}, int, error) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return NewErrorResponse(
			"InvalidDataFormatError",
			"DFE_01",
			"Invalid data format",
		), http.StatusInternalServerError, fmt.Errorf("error Marshal" + err.Error())
	}

	var request ResetPasswordVerifyCodeRequest
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

	// Generate a 6-digit OTP code
	code := fmt.Sprintf("%04d", rand.Intn(1000000))

	fakeMode := request.Fake == is.MasterKey

	// If the "fake" key exists and matches the master key, don't send SMS and set the code to "1111"
	if fakeMode {
		code = "111111"
	}

	// Set the expiration date (e.g., 5 minutes from now)
	expDate := time.Now().Add(5 * time.Minute)

	saveReq := adapter.Request{
		Method: "create",
		Data: adapter.CreateRequest{
			Collection: "otpCodes",
			Documents: []interface{}{
				OTPCode{
					Code:      code,
					ExpiredAt: expDate,
					Phone:     request.Phone,
					Email:     request.Email,
				},
			},
		},
	}

	_, err = is.Storage.Send(saveReq)
	if err != nil {
		log.Printf("Error saving OTP code to SaiStorage: %v", err)
		return NewErrorResponse(
			"ServerError",
			"SVE_06",
			"Internal server error",
		), http.StatusInternalServerError, err
	}

	log.Printf("OTP code saved to SaiStorage successfully")

	// If not a fake request, send the OTP code via SMS
	if !fakeMode {
		if request.Phone != "" {
			err = is.sendSMS(request.Phone, "Code: "+code, request.Template, request.Variables)
			if err != nil {
				log.Printf("Error sending SMS: %v", err)
				return NewErrorResponse(
					"ServerError",
					"SVE_06",
					"Internal server error",
				), http.StatusInternalServerError, err
			}
		} else {
			err = is.sendEmail(request.Email, "Code: "+code)
			if err != nil {
				log.Printf("Error sending Email: %v", err)
				return NewErrorResponse(
					"ServerError",
					"SVE_06",
					"Internal server error",
				), http.StatusInternalServerError, err
			}
		}
	}

	return NewOkResponse("Reset password code sent")
}

func (is *InternalService) sendSMS(phone, message, template string, variables map[string]interface{}) error {
	if !is.SmsEnabled {
		return nil
	}

	smsReq := SMSRequest{
		Method: "send",
		Data: SMSData{
			Phone:     phone,
			Message:   message,
			Template:  template,
			Variables: variables,
		},
	}

	smsReqJson, _ := json.Marshal(smsReq)
	_, err := http.Post(is.SmsServiceUrl, "application/json", bytes.NewBuffer(smsReqJson))

	if err != nil {
		log.Printf("Error sending SMS: %v", err)
	} else {
		log.Printf("SMS sent successfully: %s", string(smsReqJson))
	}

	return err
}

func (is *InternalService) sendEmail(email, message string) error {
	if !is.EmailEnabled {
		return nil
	}

	emailReq := EmailRequest{
		Method: "send",
		Data: EmailData{
			Sender:    is.EmailSender,
			Recipient: email,
			Body:      message,
			Subject:   "Reset Password",
		},
	}

	emailReqJson, _ := json.Marshal(emailReq)
	_, err := http.Post(is.EmailServiceUrl, "application/json", bytes.NewBuffer(emailReqJson))

	if err != nil {
		log.Printf("Error sending Email: %v", err)
	} else {
		log.Printf("Email sent successfully: %s", string(emailReqJson))
	}

	return err
}

func (is *InternalService) removeExpiredOtpCodes() {
	req := adapter.Request{
		Method: "delete",
		Data: adapter.DeleteRequest{
			Collection: "otpCodes",
			Select: map[string]interface{}{
				"expDate": map[string]interface{}{
					"$lt": time.Now(),
				},
			},
		},
	}

	resp, err := is.Storage.Send(req)
	if err != nil {
		log.Printf("Error removing expired OTP codes: %v", err)
	} else {
		log.Printf("Expired OTP codes removed successfully: %v", resp.Status)
	}
}
