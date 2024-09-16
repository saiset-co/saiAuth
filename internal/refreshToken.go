package internal

import (
	"encoding/json"
	"fmt"
	"github.com/Limpid-LLC/go-auth/internal/entities"
	"github.com/saiset-co/sai-storage-mongo/external/adapter"
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

	req := adapter.Request{
		Method: "create",
		Data: adapter.CreateRequest{
			Collection: "refreshTokens",
			Documents:  []interface{}{token},
		},
	}

	if _, err := is.Storage.Send(req); err != nil {
		return nil, fmt.Errorf("failed to save refresh token: %v", err)
	}

	return &token, nil
}

func (is InternalService) getUserByRefreshToken(refreshToken string) (*entities.User, error) {
	req := adapter.Request{
		Method: "read",
		Data: adapter.ReadRequest{
			Collection: "refreshTokens",
			Select: map[string]interface{}{
				"refresh_token": refreshToken,
				"expired_at": map[string]int64{
					"$gt": time.Now().Unix(),
				},
			},
		},
	}

	res, err := is.Storage.Send(req)
	if err != nil {
		return nil, fmt.Errorf("failed to get refresh token: %v", err)
	}

	if len(res.Result) == 0 {
		return nil, fmt.Errorf("refresh token not found")
	}

	var tokens []RefreshToken
	tokensBytes, err := json.Marshal(res.Result)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(tokensBytes, &tokens)
	if err != nil {
		return nil, err
	}

	user, err := is.UsersRepository.GetUserByID(tokens[0].UserID)

	if err != nil {
		return nil, fmt.Errorf("failed to get user by id: %v", err)
	}

	return user, nil
}

func (is InternalService) removeExpiredRefreshTokens() {
	now := time.Now().Unix()

	req := adapter.Request{
		Method: "delete",
		Data: adapter.DeleteRequest{
			Collection: "refreshTokens",
			Select: map[string]interface{}{
				"expired_at": map[string]interface{}{
					"$lt": now,
				},
			},
		},
	}

	_, err := is.Storage.Send(req)
	if err != nil {
		fmt.Printf("failed to remove expired refresh tokens: %v\n", err)
	}
}
