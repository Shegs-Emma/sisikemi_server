package token

import "time"

type Maker interface {
	CreateToken(username string, duration time.Duration, isAdmin bool) (string, *Payload, error)

	VerifyToken(token string) (*Payload, error)
}