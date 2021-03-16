package geo

import (
	"errors"
	"gopkg.in/webdeskltd/dadata.v2"
	"strconv"
)

type GeoManager interface {
	GetCoordinates(address string) (float32, float32, error)
	GetAddress(latitude, longitude float32) (string, error)
}

type Manager struct {
	token  string
	secret string
}

func NewManager(token, secret string) (*Manager, error) {
	if token == "" {
		return nil, errors.New("emoty token")
	}

	if secret == "" {
		return nil, errors.New("emoty secret")
	}

	return &Manager{token: token, secret: secret}, nil
}
