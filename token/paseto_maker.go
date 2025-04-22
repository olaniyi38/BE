package token

import (
	"errors"
	"time"

	"github.com/aead/chacha20poly1305"
	"github.com/o1egl/paseto"
)

var ErrTokenExpired = errors.New("token has expired")

type PasetoMaker struct {
	paseto       *paseto.V2
	symmetricKey []byte
}

func NewPasetoMaker(symmetricKey string) (Maker, error) {
	if len(symmetricKey) != chacha20poly1305.KeySize {
		return nil, errors.New("invalid key size")
	}

	maker := &PasetoMaker{
		paseto:       paseto.NewV2(),
		symmetricKey: []byte(symmetricKey),
	}

	return maker, nil
}

func (maker *PasetoMaker) CreateToken(username string, duration time.Duration) (string, error) {
	payload, err := NewPayload(username, duration)
	if err != nil {
		return "", err
	}

	token, err := paseto.NewV2().Encrypt(maker.symmetricKey, payload, "")
	if err != nil {
		return "", err
	}

	return token, nil
}

func (maker *PasetoMaker) VerifyToken(signed string) (*Payload, error) {
	payload := &Payload{}

	err := paseto.NewV2().Decrypt(signed, maker.symmetricKey, &payload, nil)

	if err != nil {
		return nil, err
	}

	err = payload.Valid()
	if err != nil {
		return nil, err
	}

	return payload, nil
}
