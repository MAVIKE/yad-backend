package tests

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"

	"github.com/MAVIKE/yad-backend/internal/domain"
)

func (s *APITestSuite) TestUserGetOrderOk() {
	clientId := 1
	clientType := userType

	jwt, err := s.getJWT(clientId, clientType)
	s.NoError(err)

	req, err := http.NewRequest("GET", "/api/v1/orders/1", nil)
	if err != nil {
		s.FailNow("Failed to build request", err)
	}

	req.Header.Set("Authorization", "Bearer "+jwt)
	resp := httptest.NewRecorder()
	s.app.ServeHTTP(resp, req)

	s.Require().Equal(http.StatusOK, resp.Result().StatusCode)

	var respOrder domain.Order
	respData, err := ioutil.ReadAll(resp.Body)
	s.NoError(err)

	err = json.Unmarshal(respData, &respOrder)
	s.NoError(err)

	order := orders[0]

	s.Require().Equal(order.Id, respOrder.Id)
	s.Require().Equal(order.UserId, respOrder.UserId)
	s.Require().Equal(order.RestaurantId, respOrder.RestaurantId)
	s.Require().Equal(order.CourierId, respOrder.CourierId)
	s.Require().Equal(order.DeliveryPrice, respOrder.DeliveryPrice)
	s.Require().Equal(order.TotalPrice, respOrder.TotalPrice)
	s.Require().Equal(order.Status, respOrder.Status)
	s.Require().Equal(order.Paid, respOrder.Paid)
}

func (s *APITestSuite) TestCourierGetOrderOk() {
	clientId := 1
	clientType := userType

	jwt, err := s.getJWT(clientId, clientType)
	s.NoError(err)

	req, err := http.NewRequest("GET", "/api/v1/orders/2", nil)
	if err != nil {
		s.FailNow("Failed to build request", err)
	}

	req.Header.Set("Authorization", "Bearer "+jwt)
	resp := httptest.NewRecorder()
	s.app.ServeHTTP(resp, req)

	s.Require().Equal(http.StatusOK, resp.Result().StatusCode)

	var respOrder domain.Order
	respData, err := ioutil.ReadAll(resp.Body)
	s.NoError(err)

	err = json.Unmarshal(respData, &respOrder)
	s.NoError(err)

	order := orders[1]

	s.Require().Equal(order.Id, respOrder.Id)
	s.Require().Equal(order.UserId, respOrder.UserId)
	s.Require().Equal(order.RestaurantId, respOrder.RestaurantId)
	s.Require().Equal(order.CourierId, respOrder.CourierId)
	s.Require().Equal(order.DeliveryPrice, respOrder.DeliveryPrice)
	s.Require().Equal(order.TotalPrice, respOrder.TotalPrice)
	s.Require().Equal(order.Status, respOrder.Status)
	s.Require().Equal(order.Paid, respOrder.Paid)
}

func (s *APITestSuite) TestUserGetOrderError_Forbidden() {
	clientId := 2
	clientType := userType

	jwt, err := s.getJWT(clientId, clientType)
	s.NoError(err)

	req, err := http.NewRequest("GET", "/api/v1/orders/1", nil)
	if err != nil {
		s.FailNow("Failed to build request", err)
	}

	req.Header.Set("Authorization", "Bearer "+jwt)
	resp := httptest.NewRecorder()
	s.app.ServeHTTP(resp, req)

	s.Require().Equal(http.StatusInternalServerError, resp.Result().StatusCode)
}

func (s *APITestSuite) TestCourierGetOrderError_Forbidden() {
	clientId := 2
	clientType := userType

	jwt, err := s.getJWT(clientId, clientType)
	s.NoError(err)

	req, err := http.NewRequest("GET", "/api/v1/orders/2", nil)
	if err != nil {
		s.FailNow("Failed to build request", err)
	}

	req.Header.Set("Authorization", "Bearer "+jwt)
	resp := httptest.NewRecorder()
	s.app.ServeHTTP(resp, req)

	s.Require().Equal(http.StatusInternalServerError, resp.Result().StatusCode)
}

func (s *APITestSuite) TestUserGetOrdersOk() {
	clientId := 1
	clientType := userType

	jwt, err := s.getJWT(clientId, clientType)
	s.NoError(err)

	req, err := http.NewRequest("GET", "/api/v1/users/1/orders", nil)
	if err != nil {
		s.FailNow("Failed to build request", err)
	}

	req.Header.Set("Authorization", "Bearer "+jwt)
	resp := httptest.NewRecorder()
	s.app.ServeHTTP(resp, req)

	s.Require().Equal(http.StatusOK, resp.Result().StatusCode)

	var respOrders = []*domain.Order{}
	respData, err := ioutil.ReadAll(resp.Body)
	s.NoError(err)

	err = json.Unmarshal(respData, &respOrders)
	s.NoError(err)

	clientOrders := orders[:2]

	for i := 0; i < len(clientOrders); i++ {
		s.Require().Equal(clientOrders[i].Id, respOrders[i].Id)
		s.Require().Equal(clientOrders[i].UserId, respOrders[i].UserId)
		s.Require().Equal(clientOrders[i].RestaurantId, respOrders[i].RestaurantId)
		s.Require().Equal(clientOrders[i].CourierId, respOrders[i].CourierId)
		s.Require().Equal(clientOrders[i].DeliveryPrice, respOrders[i].DeliveryPrice)
		s.Require().Equal(clientOrders[i].TotalPrice, respOrders[i].TotalPrice)
		s.Require().Equal(clientOrders[i].Status, respOrders[i].Status)
		s.Require().Equal(clientOrders[i].Paid, respOrders[i].Paid)
	}
}

func (s *APITestSuite) TestCourierGetOrdersOk() {
	clientId := 2
	clientType := courierType

	jwt, err := s.getJWT(clientId, clientType)
	s.NoError(err)

	req, err := http.NewRequest("GET", "/api/v1/couriers/2/orders", nil)
	if err != nil {
		s.FailNow("Failed to build request", err)
	}

	req.Header.Set("Authorization", "Bearer "+jwt)
	resp := httptest.NewRecorder()
	s.app.ServeHTTP(resp, req)

	fmt.Println(resp.Body.String())
	s.Require().Equal(http.StatusOK, resp.Result().StatusCode)

	var respOrder domain.Order
	respData, err := ioutil.ReadAll(resp.Body)
	s.NoError(err)

	err = json.Unmarshal(respData, &respOrder)
	s.NoError(err)

	order := orders[2]

	s.Require().Equal(order.Id, respOrder.Id)
	s.Require().Equal(order.UserId, respOrder.UserId)
	s.Require().Equal(order.RestaurantId, respOrder.RestaurantId)
	s.Require().Equal(order.CourierId, respOrder.CourierId)
	s.Require().Equal(order.DeliveryPrice, respOrder.DeliveryPrice)
	s.Require().Equal(order.TotalPrice, respOrder.TotalPrice)
	s.Require().Equal(order.Status, respOrder.Status)
	s.Require().Equal(order.Paid, respOrder.Paid)
}

func (s *APITestSuite) TestUserGetOrdersError_Forbidden() {
	clientId := 1
	clientType := userType

	jwt, err := s.getJWT(clientId, clientType)
	s.NoError(err)

	req, err := http.NewRequest("GET", "/api/v1/users/2/orders", nil)
	if err != nil {
		s.FailNow("Failed to build request", err)
	}

	req.Header.Set("Authorization", "Bearer "+jwt)
	resp := httptest.NewRecorder()
	s.app.ServeHTTP(resp, req)

	s.Require().Equal(http.StatusInternalServerError, resp.Result().StatusCode)
}

func (s *APITestSuite) TestCourierUserGetOrdersError_Forbidden() {
	clientId := 1
	clientType := userType

	jwt, err := s.getJWT(clientId, clientType)
	s.NoError(err)

	req, err := http.NewRequest("GET", "/api/v1/couriers/1/orders", nil)
	if err != nil {
		s.FailNow("Failed to build request", err)
	}

	req.Header.Set("Authorization", "Bearer "+jwt)
	resp := httptest.NewRecorder()
	s.app.ServeHTTP(resp, req)

	s.Require().Equal(http.StatusInternalServerError, resp.Result().StatusCode)
}
