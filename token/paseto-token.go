package token

import (
	"fmt"
	"time"

	"github.com/aead/chacha20poly1305"
	"github.com/o1egl/paseto"
)

// A struct that implements a paseto token function as well as a symmetricKey that will be used to create all tokens
type PasetoGoToken struct {
	paseto       *paseto.V2
	symmetricKey []byte
	SK           string
}

func (t *PasetoGoToken) GetKey() string {
	return t.SK
}

func NewPasetoClient(symmetricKey string) (GoTokens, error) {
	if len(symmetricKey) != chacha20poly1305.KeySize {
		return nil, fmt.Errorf("error. key must be %d characters long", chacha20poly1305.KeySize)
	}

	maker := &PasetoGoToken{
		paseto:       paseto.NewV2(),
		symmetricKey: []byte(symmetricKey),
		SK:           symmetricKey,
	}

	return maker, nil
}

func (p_token_maker *PasetoGoToken) CreateToken(arg_username string, duration time.Duration) (string, error) {
	payload, err := NewTokenPayload(arg_username, duration)

	if err != nil {
		return "", err
	}

	return p_token_maker.paseto.Encrypt(p_token_maker.symmetricKey, payload, nil)
}
func (p_token_maker *PasetoGoToken) VerifyToken(token string) (*TokenPayload, error) {
	token_payload := &TokenPayload{}

	err := p_token_maker.paseto.Decrypt(token, p_token_maker.symmetricKey, token_payload, nil)
	if err != nil {
		return nil, err
	}

	err = token_payload.Valid()
	if err != nil {
		return nil, err
	}
	return token_payload, err
}
