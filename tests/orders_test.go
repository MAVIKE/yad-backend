package tests

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"

	"github.com/MAVIKE/yad-backend/internal/domain"
)

func (s *APITestSuite) TestUserGetOrderOk() {
	userId := 1
	clientType := userType

	jwt, err := s.getJWT(userId, clientType)
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

	s.Require().Equal(orders[0].Id, respOrder.Id)
	s.Require().Equal(orders[0].UserId, respOrder.UserId)
	s.Require().Equal(orders[0].RestaurantId, respOrder.RestaurantId)
	s.Require().Equal(orders[0].CourierId, respOrder.CourierId)
	s.Require().Equal(orders[0].DeliveryPrice, respOrder.DeliveryPrice)
	s.Require().Equal(orders[0].TotalPrice, respOrder.TotalPrice)
	s.Require().Equal(orders[0].Status, respOrder.Status)
	s.Require().Equal(orders[0].Paid, respOrder.Paid)
}

func (s *APITestSuite) TestUserGetOrderForbidden() {
	userId := 2
	clientType := userType

	jwt, err := s.getJWT(userId, clientType)
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
