package auth

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"time"
)

type TokenManager interface {
	NewJWT(id int, clientType string, ttl time.Duration) (string, error)
}

type Manager struct {
	signingKey string
}

type tokenClaims struct {
	jwt.StandardClaims
	Id         int    `json:"id"`
	ClientType string `json:"client_type"`
}

func NewManager(signingKey string) (*Manager, error) {
	if signingKey == "" {
		return nil, errors.New("empty signing key")
	}

	return &Manager{signingKey: signingKey}, nil
}

func (m *Manager) NewJWT(id int, clientType string, ttl time.Duration) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &tokenClaims{
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(ttl).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
		id,
		clientType,
	})

	return token.SignedString([]byte(m.signingKey))
}
