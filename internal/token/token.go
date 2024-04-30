package token

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base32"
	"time"

	db "github.com/Hopertz/Hr-Be/internal/db/sqlc"
	"github.com/google/uuid"
)

const (
	ScopeAuthentication = "authentication"
)

func NewToken(user_id uuid.UUID, ttl time.Duration, scope string) (*db.CreateTokenParams, string, error) {

	token, tokenText, err := generateToken(user_id, ttl, scope)
	if err != nil {
		return nil, tokenText, err
	}

	return token, tokenText, err
}

func generateToken(user_id uuid.UUID, ttl time.Duration, scope string) (*db.CreateTokenParams, string, error) {

	token := &db.CreateTokenParams{
		UserID: user_id,
		Expiry: time.Now().Add(ttl),
		Scope:  scope,
	}

	randomBytes := make([]byte, 16)
	_, err := rand.Read(randomBytes)
	if err != nil {
		return nil, "", err
	}

	tokenPlaintext := base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString(randomBytes)
	hash := sha256.Sum256([]byte(tokenPlaintext))
	token.Hash = hash[:]
	return token, tokenPlaintext, nil
}
