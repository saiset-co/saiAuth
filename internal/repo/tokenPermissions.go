package repo

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/Limpid-LLC/go-auth/internal/entities"
	"github.com/saiset-co/sai-storage-mongo/external/adapter"
)

type TokenPermissionsRepository struct {
	Collection string
	Storage    *adapter.SaiStorage
}

func (repo TokenPermissionsRepository) RemoveExpiredTokens() {
	now := time.Now().Unix()

	req := adapter.Request{
		Method: "delete",
		Data: adapter.DeleteRequest{
			Collection: repo.Collection,
			Select: map[string]interface{}{
				"expired_at": map[string]interface{}{
					"$lt": now,
				},
			},
		},
	}

	_, err := repo.Storage.Send(req)
	if err != nil {
		fmt.Printf("failed to remove expired tokens: %v\n", err)
	}
}

func (repo TokenPermissionsRepository) FindTokenPermissions(token string, microservice string, method string) ([]entities.TokenPermission, error) {
	req := adapter.Request{
		Method: "read",
		Data: adapter.ReadRequest{
			Collection: repo.Collection,
			Select: map[string]interface{}{
				"token":                   token,
				"permission_microservice": microservice,
				"permission_method":       method,
			},
		},
	}

	res, err := repo.Storage.Send(req)
	if err != nil {
		return nil, fmt.Errorf("failed to get token permissions: %v", err)
	}

	// Convert result to slice of TokenPermission
	var tokenPermissions []entities.TokenPermission
	itemBytes, err := json.Marshal(res.Result)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(itemBytes, &tokenPermissions)
	if err != nil {
		return nil, err
	}

	return tokenPermissions, nil
}

func (repo TokenPermissionsRepository) SaveTokenPermissions(tokenPermissions []interface{}) error {
	saveReq := adapter.Request{
		Method: "create",
		Data: adapter.CreateRequest{
			Collection: repo.Collection,
			Documents:  tokenPermissions,
		},
	}

	_, err := repo.Storage.Send(saveReq)
	if err != nil {
		return fmt.Errorf("failed to save token permission: %v", err)
	}

	return nil
}

func (repo TokenPermissionsRepository) RemoveTokenPermissionsByRoleInternalID(roleInternalID string) error {
	req := adapter.Request{
		Method: "delete",
		Data: adapter.DeleteRequest{
			Collection: repo.Collection,
			Select: map[string]interface{}{
				"role_internal_id": roleInternalID,
			},
		},
	}

	_, err := repo.Storage.Send(req)
	if err != nil {
		return fmt.Errorf("failed to remove token permissions: %v", err)
	}

	return nil
}
