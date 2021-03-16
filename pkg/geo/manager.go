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

func (m *Manager) GetCoordinates(address string) (float32, float32, error) {
	d := dadata.NewDaData(m.token, m.secret)

	res, err := d.CleanAddresses(address)
	if err != nil {
		return 0, 0, err
	}

	if len(res) == 0 {
		return 0, 0, errors.New("empty result by address")
	}

	coords := res[0]

	lat, err := strconv.ParseFloat(coords.GeoLat, 32)
	if err != nil {
		return 0, 0, err
	}

	lon, err := strconv.ParseFloat(coords.GeoLon, 32)
	if err != nil {
		return 0, 0, err
	}

	return float32(lat), float32(lon), nil
}

func (m *Manager) GetAddress(latitude, longitude float32) (string, error) {
	d := dadata.NewDaData(m.token, m.secret)

	request := dadata.GeolocateRequest{Lat: latitude, Lon: longitude}
	res, err := d.GeolocateAddress(request)
	if err != nil {
		return "", err
	}

	address := res[0]
	return address.Value, nil
}

