package repo

import (
	"encoding/json"
	"fmt"
	"github.com/Limpid-LLC/go-auth/internal/entities"
	"github.com/saiset-co/sai-storage-mongo/external/adapter"
)

type UsersRepository struct {
	Collection string
	Storage    *adapter.SaiStorage
}

func (repo *UsersRepository) CreateUser(user *entities.User) error {
	req := adapter.Request{
		Method: "create",
		Data: adapter.CreateRequest{
			Collection: repo.Collection,
			Documents:  []interface{}{user},
		},
	}

	_, err := repo.Storage.Send(req)
	if err != nil {
		return fmt.Errorf("failed to create user: %v", err)
	}

	return nil
}

// UpdateUser updates a user by its internal id.
func (repo *UsersRepository) UpdateUser(user *entities.User) error {
	req := adapter.Request{
		Method: "update",
		Data: adapter.UpdateRequest{
			Collection: repo.Collection,
			Select: map[string]interface{}{
				"internal_id": user.InternalId,
			},
			Document: map[string]interface{}{"$set": user},
		},
	}

	_, err := repo.Storage.Send(req)
	if err != nil {
		return fmt.Errorf("failed to update user: %v", err)
	}

	return nil
}

func (repo *UsersRepository) GetUsersByRole(roleID string) ([]entities.User, error) {
	req := adapter.Request{
		Method: "read",
		Data: adapter.ReadRequest{
			Collection: repo.Collection,
			Select: map[string]interface{}{
				"___roles.internal_id": roleID,
			},
		},
	}

	res, err := repo.Storage.Send(req)
	if err != nil {
		return nil, err
	}

	var users []entities.User
	rByres, err := json.Marshal(res.Result)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(rByres, &users)
	if err != nil {
		return nil, err
	}

	return users, nil
}

func (repo *UsersRepository) GetUserByID(id string) (*entities.User, error) {
	req := adapter.Request{
		Method: "read",
		Data: adapter.ReadRequest{
			Collection: repo.Collection,
			Select: map[string]interface{}{
				"internal_id": id,
			},
		},
	}

	res, err := repo.Storage.Send(req)
	if err != nil || len(res.Result) == 0 {
		return nil, fmt.Errorf("user not found")
	}

	var users []entities.User
	rByres, err := json.Marshal(res.Result)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(rByres, &users)
	if err != nil {
		return nil, err
	}

	return &users[0], nil
}

func (repo *UsersRepository) GetUserByLoginAndPassword(login, hashedPassword string) (*entities.User, error) {
	req := adapter.Request{
		Method: "read",
		Data: adapter.ReadRequest{
			Collection: repo.Collection,
			Select: map[string]interface{}{
				"$or": []map[string]string{
					{"email": login, "___password": hashedPassword},
					{"phone": login, "___password": hashedPassword},
				},
			},
		},
	}

	res, err := repo.Storage.Send(req)
	if err != nil || len(res.Result) == 0 {
		return nil, fmt.Errorf("user not found")
	}

	var users []entities.User
	rByres, err := json.Marshal(res.Result)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(rByres, &users)
	if err != nil {
		return nil, err
	}

	return &users[0], nil
}

func (repo *UsersRepository) GetUserByPhone(phone string) (*entities.User, error) {
	req := adapter.Request{
		Method: "read",
		Data: adapter.ReadRequest{
			Collection: repo.Collection,
			Select: map[string]interface{}{
				"phone": phone,
			},
		},
	}

	res, err := repo.Storage.Send(req)
	if err != nil || len(res.Result) == 0 {
		return nil, fmt.Errorf("user not found")
	}

	var users []entities.User
	rByres, err := json.Marshal(res.Result)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(rByres, &users)
	if err != nil {
		return nil, err
	}

	return &users[0], nil
}

func (repo *UsersRepository) GetUserByPhoneOrEmail(phone, email string) (*entities.User, error) {
	req := adapter.Request{
		Method: "read",
		Data: adapter.ReadRequest{
			Collection: repo.Collection,
			Select: map[string]interface{}{
				"$or": []map[string]string{
					{"email": email},
					{"phone": phone},
				},
			},
		},
	}

	res, err := repo.Storage.Send(req)
	if err != nil || len(res.Result) == 0 {
		return nil, fmt.Errorf("user not found")
	}

	var users []entities.User
	rByres, err := json.Marshal(res.Result)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(rByres, &users)
	if err != nil {
		return nil, err
	}

	return &users[0], nil
}
