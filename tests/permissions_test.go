package tests

import (
	"net/http"
	"testing"
)

func TestRoleWithPermissions(t *testing.T) {
	// create role
	_, statusCode, err := sendRequest("POST", url, map[string]interface{}{
		"method": "create_role",
		"data": map[string]interface{}{
			"type":   "sto",
			"sto_id": sto_id,
			"permissions": []interface{}{
				map[string]interface{}{
					"microservice": "go-auth",
					"method":       "test_cred",
					"required_params": []interface{}{
						map[string]interface{}{
							"param":  "sto_id",
							"values": []interface{}{sto_id},
							"all":    false,
						},
						map[string]interface{}{
							"param":  "internal_id",
							"values": []interface{}{"$.internal_id"},
							"all":    false,
						},
						map[string]interface{}{
							"param":  "first_name",
							"values": []interface{}{"$.data.first_name"},
							"all":    false,
						},
					},
					"restricted_params": []interface{}{
						map[string]interface{}{
							"param":  "param1",
							"values": nil,
							"all":    true,
						},
						map[string]interface{}{
							"param":  "param2",
							"values": []interface{}{"test"},
							"all":    false,
						},
					},
				},
			},
			"data": map[string]interface{}{
				"alias":           "super-admin-2",
				"name":            "Super Admin-2",
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

	TestUserCreation(t)

	AttachRoleToUser(t)

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

	userId, err := getUserByPhone(phoneNumber)
	if err != nil {
		t.Fatalf("User not exists")
	}

	// test permissions
	_, statusCode, err = sendRequest("POST", url, map[string]interface{}{
		"method": "test_cred",
		"data": map[string]interface{}{
			"sto_id":      sto_id,
			"param2":      "test2",
			"internal_id": userId,
			"first_name":  firstName,
			"token":       token,
		},
	})

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if statusCode != http.StatusOK {
		t.Fatalf("Expected status OK, got %v", statusCode)
	}

	// test permissions (denied) - 1st case
	_, statusCode, err = sendRequest("POST", url, map[string]interface{}{
		"method": "test_cred",
		"data": map[string]interface{}{
			"token": token,
		},
	})

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if statusCode == http.StatusOK {
		t.Fatalf("Expected status not OK, got %v", statusCode)
	}

	// test permissions (denied) - 2nd case
	_, statusCode, err = sendRequest("POST", url, map[string]interface{}{
		"method": "test_cred",
		"data": map[string]interface{}{
			"sto_id": sto_id,
			"param1": "some_value",
			"token":  token,
		},
	})

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if statusCode == http.StatusOK {
		t.Fatalf("Expected status not OK, got %v", statusCode)
	}

	// test permissions (denied) - 4rd case
	_, statusCode, err = sendRequest("POST", url, map[string]interface{}{
		"method": "test_cred",
		"data": map[string]interface{}{
			"sto_id": sto_id,
			"param2": "test",
			"token":  token,
		},
	})

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if statusCode == http.StatusOK {
		t.Fatalf("Expected status not OK, got %v", statusCode)
	}

	// test permissions (denied) - 5th case
	_, statusCode, err = sendRequest("POST", url, map[string]interface{}{
		"method": "test_cred",
		"data": map[string]interface{}{
			"sto_id": sto_id,
			"param2": "test",
		},
	})

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if statusCode == http.StatusOK {
		t.Fatalf("Expected status not OK, got %v", statusCode)
	}
}

func AttachRoleToUser(t *testing.T) {
	roleId, err := getRoleByStoId(sto_id, "super-admin-2")

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
}
