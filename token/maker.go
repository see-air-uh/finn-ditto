package token

import (
	"time"
)

type GoTokens interface {
	CreateToken(arg_username string, duration time.Duration) (string, error)
	VerifyToken(token string) (*TokenPayload, error)
}
