package token

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

var jwtSigningMethod = jwt.SigningMethodHS256

const (
	ClaimExpiry   = "exp"
	ClaimIssuedAt = "iat"
)

type JWTMaker struct {
	secretKey string
}

func NewJWTMaker(secretKey string) (Maker, error) {
	if len(secretKey) < 8 {
		return nil, fmt.Errorf("secret key must be at least 8 characters")
	}
	return JWTMaker{secretKey: secretKey}, nil
}

func (maker JWTMaker) CreateToken(username string, duration time.Duration) (string, error) {
	payload, err := NewPayload(username, duration)
	if err != nil {
		return "", err
	}

	claims := jwt.MapClaims{
		"id":          payload.ID,
		"username":    payload.Username,
		ClaimExpiry:   timeToUnix(payload.ExpiredAt),
		ClaimIssuedAt: timeToUnix(payload.IssuedAt),
	}

	fmt.Println(claims)

	token := jwt.NewWithClaims(jwtSigningMethod, claims)

	signed, err := token.SignedString([]byte(maker.secretKey))
	if err != nil {
		return "", err
	}

	return signed, nil
}

func (maker JWTMaker) VerifyToken(signedToken string) (*Payload, error) {

	//validate signing method
	parseOption := jwt.WithValidMethods([]string{jwtSigningMethod.Name})

	token, err := jwt.Parse(signedToken, func(t *jwt.Token) (interface{}, error) {
		// check expiration
		tokenExpired, err := checkTokenExpiration(t)
		if err != nil {
			return nil, err
		}

		if tokenExpired {
			return nil, jwt.ErrTokenExpired
		}
		return []byte(maker.secretKey), nil
	}, parseOption)

	if err != nil {
		return nil, err
	}

	mapClaims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("invalid claims")
	}

	username, ok := mapClaims["username"].(string)
	if !ok {
		return nil, fmt.Errorf("username not found")
	}

	id, ok := mapClaims["id"].(string)

	if !ok {
		return nil, fmt.Errorf("token id not found")
	}

	parsedID, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid uuid found %v", err)
	}

	payload := &Payload{
		ID:       parsedID,
		Username: username,
	}

	return payload, nil
}

func checkTokenExpiration(t *jwt.Token) (bool, error) {
	expirationTime, err := t.Claims.GetExpirationTime()
	if err != nil {
		return false, err
	}

	now := time.Now()
	expired := now.After(expirationTime.Time)
	return expired, nil
}

func timeToUnix(t time.Time) float64 {
	return float64(t.Unix())
}
