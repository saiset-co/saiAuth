package repo

import (
	"encoding/json"
	"fmt"
	"github.com/Limpid-LLC/go-auth/internal/entities"
	"github.com/Limpid-LLC/go-auth/internal/storage"
	"time"
)

type TokenPermissionsRepository struct {
	Collection string
	Storage    *storage.SaiStorage
}

func (repo TokenPermissionsRepository) RemoveExpiredTokens() {
	now := time.Now().Unix()

	req := storage.SaiStorageRemoveRequest{
		Collection: repo.Collection,
		Select: map[string]interface{}{
			"expired_at": map[string]interface{}{
				"$lt": now,
			},
		},
	}

	_, err := repo.Storage.Remove(req)
	if err != nil {
		fmt.Printf("failed to remove expired tokens: %v\n", err)
	}
}

func (repo TokenPermissionsRepository) FindTokenPermissions(token string, microservice string, method string) ([]entities.TokenPermission, error) {

	req := storage.SaiStorageGetRequest{
		Collection: repo.Collection,
		Select: map[string]interface{}{
			"token":                   token,
			"permission_microservice": microservice,
			"permission_method":       method,
		},
	}

	res, err := repo.Storage.GetEncoded(req)
	if err != nil {
		return nil, fmt.Errorf("failed to get token permissions: %v", err)
	}

	// Convert result to slice of TokenPermission
	var tokenPermissions []entities.TokenPermission
	for _, item := range res.Result {
		var tp entities.TokenPermission
		err = json.Unmarshal(item, &tp)
		if err != nil {
			return nil, err
		}
		tokenPermissions = append(tokenPermissions, tp)
	}

	return tokenPermissions, nil
}

func (repo TokenPermissionsRepository) SaveTokenPermissions(tokenPermissions []entities.TokenPermission) error {

	for _, tp := range tokenPermissions {
		saveReq := storage.SaiStorageSaveRequest{
			Collection: repo.Collection,
			Data:       tp,
		}
		_, err := repo.Storage.Save(saveReq)
		if err != nil {
			return fmt.Errorf("failed to save token permission: %v", err)
		}
	}
	return nil
}

func (repo TokenPermissionsRepository) RemoveTokenPermissionsByRoleInternalID(roleInternalID string) error {
	req := storage.SaiStorageRemoveRequest{
		Collection: repo.Collection,
		Select: map[string]interface{}{
			"role_internal_id": roleInternalID,
		},
	}

	_, err := repo.Storage.Remove(req)
	if err != nil {
		return fmt.Errorf("failed to remove token permissions: %v", err)
	}

	return nil
}
