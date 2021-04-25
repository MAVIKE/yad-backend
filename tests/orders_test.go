package tests

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"

	"github.com/MAVIKE/yad-backend/internal/domain"
)

func (s *APITestSuite) TestUserGetOrder() {
	order := domain.Order{
		Id:            1,
		UserId:        1,
		RestaurantId:  1,
		CourierId:     0,
		DeliveryPrice: 100,
		TotalPrice:    900,
		Status:        0,
		Paid:          nil,
	}

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

	fmt.Println(resp.Body.String())
	s.Require().Equal(http.StatusOK, resp.Result().StatusCode)

	var respOrder domain.Order
	respData, err := ioutil.ReadAll(resp.Body)
	s.NoError(err)

	err = json.Unmarshal(respData, &respOrder)
	s.NoError(err)

	// s.Require().Equal(1, len(respOffers.Data))
	s.Require().Equal(order.Id, respOrder.Id)
	s.Require().Equal(order.UserId, respOrder.UserId)
	s.Require().Equal(order.RestaurantId, respOrder.RestaurantId)
	s.Require().Equal(order.CourierId, respOrder.CourierId)
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

	fmt.Println(resp.Body.String())
	s.Require().Equal(http.StatusInternalServerError, resp.Result().StatusCode)
}
