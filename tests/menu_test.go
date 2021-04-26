package tests

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"

	"github.com/MAVIKE/yad-backend/internal/domain"
)

func (s *APITestSuite) TestRestaurantGetAllMenuItemsOk() {
	clientId := 1
	clientType := restaurantType
	jwt, err := s.getJWT(clientId, clientType)
	s.NoError(err)

	req, err := http.NewRequest("GET", "/api/v1/restaurants/1/menu/", nil)
	if err != nil {
		s.FailNow("Failed to build request", err)
	}

	req.Header.Set("Authorization", "Bearer "+jwt)
	resp := httptest.NewRecorder()
	s.app.ServeHTTP(resp, req)

	s.Require().Equal(http.StatusOK, resp.Result().StatusCode)

	var respMenu []*domain.MenuItem
	respData, err := ioutil.ReadAll(resp.Body)
	s.NoError(err)

	err = json.Unmarshal(respData, &respMenu)
	s.NoError(err)

	s.Require().Equal(len(menuItems), len(respMenu))
	for i := 0; i < len(menuItems); i++ {
		s.Require().Equal(menuItems[i].Id, respMenu[i].Id)
		s.Require().Equal(menuItems[i].Image, respMenu[i].Image)
		s.Require().Equal(menuItems[i].RestaurantId, respMenu[i].RestaurantId)
		s.Require().Equal(menuItems[i].Price, respMenu[i].Price)
		s.Require().Equal(menuItems[i].Description, respMenu[i].Description)
		s.Require().Equal(menuItems[i].Title, respMenu[i].Title)
	}
}

func (s *APITestSuite) TestRestaurantGetAllMenuItemsForbidden() {
	clientId := 1
	clientType := courierType
	jwt, err := s.getJWT(clientId, clientType)
	s.NoError(err)

	req, err := http.NewRequest("GET", "/api/v1/restaurants/1/menu/", nil)
	if err != nil {
		s.FailNow("Failed to build request", err)
	}

	req.Header.Set("Authorization", "Bearer "+jwt)
	resp := httptest.NewRecorder()
	s.app.ServeHTTP(resp, req)

	s.Require().Equal(http.StatusInternalServerError, resp.Result().StatusCode)
}


func (s *APITestSuite) TestRestaurantGetMenuItemByIdOk() {
	clientId := 1
	clientType := userType
	jwt, err := s.getJWT(clientId, clientType)
	s.NoError(err)

	req, err := http.NewRequest("GET", "/api/v1/restaurants/1/menu/1", nil)
	if err != nil {
		s.FailNow("Failed to build request", err)
	}

	req.Header.Set("Authorization", "Bearer "+jwt)
	resp := httptest.NewRecorder()
	s.app.ServeHTTP(resp, req)

	s.Require().Equal(http.StatusOK, resp.Result().StatusCode)

	var respMenu domain.MenuItem
	respData, err := ioutil.ReadAll(resp.Body)
	s.NoError(err)

	err = json.Unmarshal(respData, &respMenu)
	s.NoError(err)

	s.Require().Equal(menuItems[0].Id, respMenu.Id)
	s.Require().Equal(menuItems[0].Image, respMenu.Image)
	s.Require().Equal(menuItems[0].RestaurantId, respMenu.RestaurantId)
	s.Require().Equal(menuItems[0].Price, respMenu.Price)
	s.Require().Equal(menuItems[0].Description, respMenu.Description)
	s.Require().Equal(menuItems[0].Title, respMenu.Title)
}

func (s *APITestSuite) TestRestaurantGetMenuItemByIdForbidden() {
	clientId := 2
	clientType := restaurantType
	jwt, err := s.getJWT(clientId, clientType)
	s.NoError(err)

	req, err := http.NewRequest("GET", "/api/v1/restaurants/1/menu/1", nil)
	if err != nil {
		s.FailNow("Failed to build request", err)
	}

	req.Header.Set("Authorization", "Bearer "+jwt)
	resp := httptest.NewRecorder()
	s.app.ServeHTTP(resp, req)

	s.Require().Equal(http.StatusInternalServerError, resp.Result().StatusCode)
}


func (s *APITestSuite) TestRestaurantGetAllCategoryMenuItemsOk() {
	clientId := 1
	clientType := userType
	jwt, err := s.getJWT(clientId, clientType)
	s.NoError(err)

	req, err := http.NewRequest("GET", "/api/v1/restaurants/1/categories/1/menu", nil)
	if err != nil {
		s.FailNow("Failed to build request", err)
	}

	req.Header.Set("Authorization", "Bearer "+jwt)
	resp := httptest.NewRecorder()
	s.app.ServeHTTP(resp, req)

	s.Require().Equal(http.StatusOK, resp.Result().StatusCode)

	var respMenu []*domain.MenuItem
	respData, err := ioutil.ReadAll(resp.Body)
	s.NoError(err)

	err = json.Unmarshal(respData, &respMenu)
	s.NoError(err)

	s.Require().Equal(2, len(respMenu))
	for i := 0; i < len(respMenu); i++ {
		s.Require().Equal(menuItems[i].Id, respMenu[i].Id)
		s.Require().Equal(menuItems[i].Image, respMenu[i].Image)
		s.Require().Equal(menuItems[i].RestaurantId, respMenu[i].RestaurantId)
		s.Require().Equal(menuItems[i].Price, respMenu[i].Price)
		s.Require().Equal(menuItems[i].Description, respMenu[i].Description)
		s.Require().Equal(menuItems[i].Title, respMenu[i].Title)
	}
}

func (s *APITestSuite) TestRestaurantGetAllCategoryItemsForbidden() {
	clientId := 1
	clientType := courierType
	jwt, err := s.getJWT(clientId, clientType)
	s.NoError(err)

	req, err := http.NewRequest("GET", "/api/v1/restaurants/1/categories/1/menu", nil)
	if err != nil {
		s.FailNow("Failed to build request", err)
	}

	req.Header.Set("Authorization", "Bearer "+jwt)
	resp := httptest.NewRecorder()
	s.app.ServeHTTP(resp, req)

	s.Require().Equal(http.StatusInternalServerError, resp.Result().StatusCode)
}
