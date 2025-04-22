package token

import "time"

// Maker is an interface for managing tokens
type Maker interface {
	//creates a token with the username as a payload
	CreateToken(username string, duration time.Duration) (string, error)

	//used to verify a token then return its payload
	VerifyToken(token string) (payload *Payload, err error)
}
