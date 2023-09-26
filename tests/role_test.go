package tests

import (
	"fmt"
	"net/http"
	"testing"
)

func TestCreateRole(t *testing.T) {
	// create role
	_, statusCode, err := sendRequest("POST", url, map[string]interface{}{
		"method": "create_role",
		"data": map[string]interface{}{
			"type":   "sto",
			"sto_id": sto_id,
			"permissions": []interface{}{
				map[string]interface{}{
					"microservice":      "go-auth",
					"method":            "test",
					"required_params":   []interface{}{},
					"restricted_params": []interface{}{},
				},
			},
			"data": map[string]interface{}{
				"alias":           "super-admin",
				"name":            "Super Admin",
				"any_other_field": "blabla",
			},
		},
	})
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if statusCode != http.StatusOK {
		t.Fatalf("Expected status OK, got %v", statusCode)
	}
}

func getRoleByStoId(stoID, roleAlias string) (string, error) {
	roleRequestData := map[string]interface{}{
		"method": "get_roles",
		"data": map[string]interface{}{
			"sto_id":     stoID,
			"data.alias": roleAlias,
		},
	}

	resp, status, _ := sendRequest("POST", url, roleRequestData)
	if status != http.StatusOK {
		return "", fmt.Errorf("expected status OK, got %v", status)
	}

	return resp.(map[string]interface{})["result"].([]interface{})[0].(map[string]interface{})["internal_id"].(string), nil
}

func getUserByPhone(phone string) (string, error) {
	roleRequestData := map[string]interface{}{
		"method": "get_users",
		"data": map[string]interface{}{
			"phone": phone,
		},
	}

	resp, status, _ := sendRequest("POST", url, roleRequestData)
	if status != http.StatusOK {
		return "", fmt.Errorf("expected status OK, got %v", status)
	}

	return resp.(map[string]interface{})["result"].([]interface{})[0].(map[string]interface{})["internal_id"].(string), nil
}

func TestRoleCreationAttachDetach(t *testing.T) {
	TestUserCreation(t)
	TestCreateRole(t)

	roleId, err := getRoleByStoId(sto_id, "super-admin")

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	userId, err := getUserByPhone(phoneNumber)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// attach role to user
	_, statusCode, err := sendRequest("POST", url, map[string]interface{}{
		"method": "attach_role",
		"data": map[string]interface{}{
			"role_id": roleId,
			"user_id": userId,
		},
	})

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if statusCode != http.StatusOK {
		t.Fatalf("Expected status OK, got %v", statusCode)
	}

	// login
	resp, statusCode, err := sendRequest("POST", url, map[string]interface{}{
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

	resultMap := resp.(map[string]interface{})["result"]
	accessTokens := resultMap.(map[string]interface{})["accessTokens"].([]interface{})
	if len(accessTokens) != 1 {
		t.Fatalf("Expected 1 access token, got %v", len(accessTokens))
	}
	token := accessTokens[0].(map[string]interface{})["token"].(string)

	if token == "" {
		t.Fatalf("Expected token, got %v", token)
	}

	// detach role from user
	_, statusCode, err = sendRequest("POST", url, map[string]interface{}{
		"method": "detach_role",
		"data": map[string]interface{}{
			"role_id": roleId,
			"user_id": userId,
		},
	})

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if statusCode != http.StatusOK {
		t.Fatalf("Expected status OK, got %v", statusCode)
	}

	// login
	resp, statusCode, err = sendRequest("POST", url, map[string]interface{}{
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

	resultMap = resp.(map[string]interface{})["result"]
	accessTokens = resultMap.(map[string]interface{})["accessTokens"].([]interface{})
	if len(accessTokens) != 0 {
		t.Fatalf("Expected 0 access token, got %v", len(accessTokens))
	}

	// attach role to user
	_, statusCode, err = sendRequest("POST", url, map[string]interface{}{
		"method": "attach_role",
		"data": map[string]interface{}{
			"role_id": roleId,
			"user_id": userId,
		},
	})

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if statusCode != http.StatusOK {
		t.Fatalf("Expected status OK, got %v", statusCode)
	}

	// login
	resp, statusCode, err = sendRequest("POST", url, map[string]interface{}{
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

	// check in response exist access token
	resultMap = resp.(map[string]interface{})["result"]
	accessTokens = resultMap.(map[string]interface{})["accessTokens"].([]interface{})
	if len(accessTokens) != 1 {
		t.Fatalf("Expected 1 access token, got %v", len(accessTokens))
	}

	// delete role
	_, statusCode, err = sendRequest("POST", url, map[string]interface{}{
		"method": "delete_role",
		"data": map[string]interface{}{
			"internal_id": roleId,
		},
	})

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if statusCode != http.StatusOK {
		t.Fatalf("Expected status OK, got %v", statusCode)
	}

	// login
	resp, statusCode, err = sendRequest("POST", url, map[string]interface{}{
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

	// check in response not exists access token
	resultMap = resp.(map[string]interface{})["result"]
	accessTokens = resultMap.(map[string]interface{})["accessTokens"].([]interface{})
	if len(accessTokens) != 0 {
		t.Fatalf("Expected 0 access token, got %v", len(accessTokens))
	}
}
