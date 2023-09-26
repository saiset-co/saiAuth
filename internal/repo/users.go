package repo

import (
	"encoding/json"
	"fmt"
	"github.com/Limpid-LLC/go-auth/internal/entities"
	"github.com/Limpid-LLC/go-auth/internal/storage"
)

type UsersRepository struct {
	Collection string
	Storage    *storage.SaiStorage
}

func (repo *UsersRepository) CreateUser(user *entities.User) error {
	_, err := repo.Storage.Save(storage.SaiStorageSaveRequest{
		Collection: repo.Collection,
		Data:       user,
	})
	if err != nil {
		return fmt.Errorf("failed to create user: %v", err)
	}

	return nil
}

// UpdateUser updates a user by its internal id.
func (repo *UsersRepository) UpdateUser(user *entities.User) error {
	_, err := repo.Storage.Update(storage.SaiStorageUpdateRequest{
		Collection: repo.Collection,
		Select: map[string]interface{}{
			"internal_id": user.InternalId,
		},
		Data: &user,
	})

	if err != nil {
		return fmt.Errorf("failed to update user: %v", err)
	}

	return nil
}

func (repo *UsersRepository) GetUsersByRole(roleID string) ([]entities.User, error) {
	// Fetch users who have this role
	req := storage.SaiStorageGetRequest{
		Collection: repo.Collection,
		Select: map[string]interface{}{
			"___roles.internal_id": roleID,
		},
	}

	res, err := repo.Storage.GetEncoded(req)
	if err != nil {
		return nil, err
	}

	var users []entities.User
	// For each user, decode the user object
	for _, r := range res.Result {
		var user entities.User
		err = json.Unmarshal(r, &user)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}

func (repo *UsersRepository) GetUserByID(id string) (*entities.User, error) {
	res, err := repo.Storage.GetEncoded(storage.SaiStorageGetRequest{
		Collection: repo.Collection,
		Select: map[string]interface{}{
			"internal_id": id,
		},
	})
	if err != nil || len(res.Result) == 0 {
		return nil, fmt.Errorf("user not found")
	}

	var user entities.User
	if err := json.Unmarshal(res.Result[0], &user); err != nil {
		return nil, fmt.Errorf("failed to unmarshal user: %v", err)
	}

	return &user, nil
}

func (repo *UsersRepository) GetUserByLoginAndPassword(login, hashedPassword string) (*entities.User, error) {
	res, err := repo.Storage.GetEncoded(storage.SaiStorageGetRequest{
		Collection: repo.Collection,
		Select: map[string][]map[string]string{
			"$or": {
				{"email": login, "___password": hashedPassword},
				{"phone": login, "___password": hashedPassword},
			},
		},
	})
	if err != nil || len(res.Result) == 0 {
		return nil, fmt.Errorf("user not found or password repo incorrect")
	}

	var user entities.User
	if err := json.Unmarshal(res.Result[0], &user); err != nil {
		return nil, fmt.Errorf("failed to unmarshal user: %v", err)
	}

	return &user, nil
}

func (repo *UsersRepository) GetUserByPhone(phone string) (*entities.User, error) {
	res, err := repo.Storage.GetEncoded(storage.SaiStorageGetRequest{
		Collection: repo.Collection,
		Select: map[string]string{
			"phone": phone,
		},
	})
	if err != nil || len(res.Result) == 0 {
		return nil, fmt.Errorf("user not found")
	}

	var user entities.User
	if err := json.Unmarshal(res.Result[0], &user); err != nil {
		return nil, fmt.Errorf("failed to unmarshal user: %v", err)
	}

	return &user, nil
}

func (repo *UsersRepository) GetUserByPhoneOrEmail(phone, email string) (*entities.User, error) {
	res, err := repo.Storage.GetEncoded(storage.SaiStorageGetRequest{
		Collection: repo.Collection,
		Select: map[string]interface{}{
			"$or": []map[string]string{
				{"email": email},
				{"phone": phone},
			},
		},
	})
	if err != nil || len(res.Result) == 0 {
		return nil, fmt.Errorf("user not found")
	}

	var user entities.User
	if err := json.Unmarshal(res.Result[0], &user); err != nil {
		return nil, fmt.Errorf("failed to unmarshal user: %v", err)
	}

	return &user, nil
}
