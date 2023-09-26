package internal

import (
	"encoding/json"
	"fmt"
	"github.com/Limpid-LLC/go-auth/internal/entities"
	"github.com/Limpid-LLC/go-auth/internal/storage"
	"time"
)

type RefreshToken struct {
	RefreshToken string `json:"refresh_token"`
	ExpiredAt    int64  `json:"expired_at"`
	UserID       string `json:"user_id"`
}

func (is *InternalService) generateRefreshToken(user *entities.User) (*RefreshToken, error) {
	// Generate a random refresh token
	refreshToken, err := generateRandomToken(64)

	if err != nil {
		return nil, fmt.Errorf("failed to generate refresh token: %v", err)
	}

	// Define the expiration time for the refresh token
	expiredAt := time.Now().Add(is.TokenExpirations.RefreshToken).Unix()

	token := RefreshToken{
		RefreshToken: refreshToken,
		ExpiredAt:    expiredAt,
		UserID:       user.InternalId,
	}

	// Store the refresh token in the database
	req := storage.SaiStorageSaveRequest{
		Collection: "refreshTokens",
		Data:       token,
	}

	if _, err := is.Storage.Save(req); err != nil {
		return nil, fmt.Errorf("failed to save refresh token: %v", err)
	}

	return &token, nil
}

func (is InternalService) getUserByRefreshToken(refreshToken string) (*entities.User, error) {

	req := storage.SaiStorageGetRequest{
		Collection: "refreshTokens",
		Select: map[string]interface{}{
			"refresh_token": refreshToken,
			"expired_at": map[string]int64{
				"$gt": time.Now().Unix(),
			},
		},
	}

	res, err := is.Storage.GetEncoded(req)
	if err != nil {
		return nil, fmt.Errorf("failed to get refresh token: %v", err)
	}

	if len(res.Result) == 0 {
		return nil, fmt.Errorf("refresh token not found")
	}

	var token RefreshToken
	err = json.Unmarshal(res.Result[0], &token)
	if err != nil {
		return nil, err
	}

	user, err := is.UsersRepository.GetUserByID(token.UserID)

	if err != nil {
		return nil, fmt.Errorf("failed to get user by id: %v", err)
	}

	return user, nil
}

func (is InternalService) removeExpiredRefreshTokens() {
	now := time.Now().Unix()

	req := storage.SaiStorageRemoveRequest{
		Collection: "refreshTokens",
		Select: map[string]interface{}{
			"expired_at": map[string]interface{}{
				"$lt": now,
			},
		},
	}

	_, err := is.Storage.Remove(req)
	if err != nil {
		fmt.Printf("failed to remove expired refresh tokens: %v\n", err)
	}
}
