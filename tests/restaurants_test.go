package tests

import (
	"encoding/json"
	"github.com/MAVIKE/yad-backend/internal/domain"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
)

func (s *APITestSuite) TestGetAllRestaurantsOk() {
	userId := 1
	clientType := userType

	jwt, err := s.getJWT(userId, clientType)
	s.NoError(err)

	req, err := http.NewRequest("GET", "/api/v1/restaurants/", nil)
	if err != nil {
		s.FailNow("Failed to build request", err)
	}

	req.Header.Set("Authorization", "Bearer "+jwt)
	resp := httptest.NewRecorder()
	s.app.ServeHTTP(resp, req)

	s.Require().Equal(http.StatusOK, resp.Result().StatusCode)

	respData, err := ioutil.ReadAll(resp.Body)
	s.NoError(err)

	var respRestaurant []*domain.Restaurant
	err = json.Unmarshal(respData, &respRestaurant)
	s.NoError(err)

	s.Require().Equal(len(restaurants), len(respRestaurant))
	for i := 0; i < len(respRestaurant); i++ {
		s.Require().Equal(restaurants[i].Id, respRestaurant[i].Id)
		s.Require().Equal(restaurants[i].Image, respRestaurant[i].Image)
		s.Require().Equal(restaurants[i].WorkingStatus, respRestaurant[i].WorkingStatus)
		s.Require().Equal(restaurants[i].Address, respRestaurant[i].Address)
		s.Require().Equal(restaurants[i].Phone, respRestaurant[i].Phone)
		s.Require().Equal(restaurants[i].Name, respRestaurant[i].Name)
	}
}

func (s *APITestSuite) TestGetAllRestaurantsForbidden() {
	userId := 1
	clientType := courierType

	jwt, err := s.getJWT(userId, clientType)
	s.NoError(err)

	req, err := http.NewRequest("GET", "/api/v1/restaurants/", nil)
	if err != nil {
		s.FailNow("Failed to build request", err)
	}

	req.Header.Set("Authorization", "Bearer "+jwt)
	resp := httptest.NewRecorder()
	s.app.ServeHTTP(resp, req)

	s.Require().Equal(http.StatusInternalServerError, resp.Result().StatusCode)
}

func (s *APITestSuite) TestGetRestaurantByIdOk() {
	userId := 1
	clientType := userType

	jwt, err := s.getJWT(userId, clientType)
	s.NoError(err)

	req, err := http.NewRequest("GET", "/api/v1/restaurants/1", nil)
	if err != nil {
		s.FailNow("Failed to build request", err)
	}

	req.Header.Set("Authorization", "Bearer "+jwt)
	resp := httptest.NewRecorder()
	s.app.ServeHTTP(resp, req)

	s.Require().Equal(http.StatusOK, resp.Result().StatusCode)

	respData, err := ioutil.ReadAll(resp.Body)
	s.NoError(err)

	var respRestaurant domain.Restaurant
	err = json.Unmarshal(respData, &respRestaurant)
	s.NoError(err)

	s.Require().Equal(restaurants[0].Id, respRestaurant.Id)
	s.Require().Equal(restaurants[0].Image, respRestaurant.Image)
	s.Require().Equal(restaurants[0].WorkingStatus, respRestaurant.WorkingStatus)
	s.Require().Equal(restaurants[0].Address, respRestaurant.Address)
	s.Require().Equal(restaurants[0].Phone, respRestaurant.Phone)
	s.Require().Equal(restaurants[0].Name, respRestaurant.Name)
}

func (s *APITestSuite) TestGetRestaurantByIdForbidden() {
	restaurantId := 2
	clientType := restaurantType

	jwt, err := s.getJWT(restaurantId, clientType)
	s.NoError(err)

	req, err := http.NewRequest("GET", "/api/v1/restaurants/1", nil)
	if err != nil {
		s.FailNow("Failed to build request", err)
	}

	req.Header.Set("Authorization", "Bearer "+jwt)
	resp := httptest.NewRecorder()
	s.app.ServeHTTP(resp, req)

	s.Require().Equal(http.StatusInternalServerError, resp.Result().StatusCode)
}
