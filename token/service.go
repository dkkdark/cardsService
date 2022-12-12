package token

import (
	"encoding/json"
	"fmt"
	jose "github.com/dvsekhvalnov/jose2go"
	rsa "github.com/dvsekhvalnov/jose2go/keys/rsa"
)

type ServiceImpl struct {
	publicKey  string
	privateKey string
}

func New(publicKey, privateKey string) *ServiceImpl {
	return &ServiceImpl{
		publicKey:  publicKey,
		privateKey: privateKey,
	}
}

func (s *ServiceImpl) GetToken(userId, role string) (string, error) {
	privateKeyRead, err := rsa.ReadPrivate([]byte(s.privateKey))

	if err != nil {
		return "", fmt.Errorf("getToken, rsa.ReadPrivate error, err: %w", err)
	}

	payload, err := json.Marshal(&Token{
		UserId: userId,
		Role:   role,
	})

	token, err := jose.Sign(string(payload), jose.RS256, privateKeyRead)

	if err != nil {
		return "", fmt.Errorf("getToken jose.Sign error, err: %w", err)
	}

	return token, nil
}

func (s *ServiceImpl) ParseToken(token string) (*Token, error) {
	publicKeyRead, err := rsa.ReadPublic([]byte(s.publicKey))
	if err != nil {
		return nil, fmt.Errorf("parseToken, rsa.ReadPublic error, err: %w", err)
	}

	payloadByte, _, err := jose.DecodeBytes(token, publicKeyRead)
	if err != nil {
		return nil, fmt.Errorf("parseToken, rsa.DecodeBytes error, err: %w", err)
	}

	t := &Token{}
	err = json.Unmarshal(payloadByte, t)
	if err != nil {
		return nil, fmt.Errorf("parseToken, json.Unmarshal error, err: %w", err)
	}
	return t, nil
}
