package internal

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"time"
)

func generateRandomToken(n int) (string, error) {
	bytes := make([]byte, n)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

func (is *InternalService) hashAndSaltPassword(password string) string {
	saltedPassword := password + is.Salt
	hashedPassword := sha256.Sum256([]byte(saltedPassword))
	return hex.EncodeToString(hashedPassword[:])
}

func startCleanupRoutine(ctx context.Context, interval time.Duration, cleanCallback func()) {
	ticker := time.NewTicker(interval)
	for {
		select {
		case <-ticker.C:
			cleanCallback()
		case <-ctx.Done():
			ticker.Stop()
			return
		}
	}
}
