package tests

import (
	"net/http"
	"testing"
)

const (
	url         = "http://localhost:9001"
	phoneNumber = "+33333333333"
	email       = "john3@example.com"
	password    = "+12345678915"
	firstName   = "John"
	secondName  = "Snow"

	sto_id = "-1"
)

func TestUserCreation(t *testing.T) {

	// remove all users with phone +33333333333
	_, statusCode, err := sendRequest("POST", url, map[string]interface{}{
		"method": "delete_users",
		"data": map[string]interface{}{
			"$or": []map[string]interface{}{
				{"phone": phoneNumber},
				{"email": email},
			},
		},
	})
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if statusCode != http.StatusOK {
		t.Fatalf("Expected status OK, got %v", statusCode)
	}

	// send otp
	_, statusCode, err = sendRequest("POST", url, map[string]interface{}{
		"method": "send_verify_code",
		"data": map[string]interface{}{
			"phone": phoneNumber,
			"fake":  "mstfiqalx",
		},
	})
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if statusCode != http.StatusOK {
		t.Fatalf("Expected status OK, got %v", statusCode)
	}

	// try create user with not valid otp
	_, statusCode, err = sendRequest("POST", url, map[string]interface{}{
		"method": "sign_up",
		"data": map[string]interface{}{
			"email":    email,
			"phone":    phoneNumber,
			"otp_code": "2222", // invalid otp code
			"password": password,
			"data": map[string]interface{}{
				"first_name": firstName,
				"last_name":  secondName,
			},
		},
	})
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if statusCode == http.StatusOK {
		t.Fatalf("Expected status not OK, got %v", statusCode)
	}

	// create user
	_, statusCode, err = sendRequest("POST", url, map[string]interface{}{
		"method": "sign_up",
		"data": map[string]interface{}{
			"email":    email,
			"phone":    phoneNumber,
			"otp_code": "1111", // valid otp code
			"password": password,
			"data": map[string]interface{}{
				"first_name": firstName,
				"last_name":  secondName,
			},
		},
	})
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if statusCode != http.StatusOK {
		t.Fatalf("Expected status OK, got %v", statusCode)
	}

	// try login with wrong password
	_, statusCode, err = sendRequest("POST", url, map[string]interface{}{
		"method": "sign_in",
		"data": map[string]interface{}{
			"login":    phoneNumber,
			"password": "wrongpassword", // wrong password
		},
	})
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if statusCode == http.StatusOK {
		t.Fatalf("Expected status not OK, got %v", statusCode)
	}

	// login
	_, statusCode, err = sendRequest("POST", url, map[string]interface{}{
		"method": "sign_in",
		"data": map[string]interface{}{
			"login":    phoneNumber,
			"password": password, // correct password
		},
	})
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if statusCode != http.StatusOK {
		t.Fatalf("Expected status OK, got %v", statusCode)
	}
}
