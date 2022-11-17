package token

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

type TokenPayload struct {
	ID        uuid.UUID `json:"id"`
	Username  string    `json:"username"`
	IssuedAt  time.Time `json:"issuedAt"`
	ExpiredAt time.Time `json:"expiredAt"`
}

// NewTokenPayload creates a new token payload with a username, the issue date, expired at date as well as a unique ID
func NewTokenPayload(username string, duration time.Duration) (*TokenPayload, error) {
	token_id, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}

	token_payload := &TokenPayload{
		ID:        token_id,
		Username:  username,
		IssuedAt:  time.Now(),
		ExpiredAt: time.Now().Add(duration),
	}
	return token_payload, nil
}

// Valid checks if the token payload is expired
func (token_payload *TokenPayload) Valid() error {
	if time.Now().After(token_payload.ExpiredAt) {
		return errors.New("error. token expired")
	}
	return nil
}
