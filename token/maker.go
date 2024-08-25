package token

import "time"

type Maker interface {
	CreateToken(username string, duration time.Duration, isAdmin bool) (string, error)

	VerifyToken(token string) (*Payload, error)
}