// +build e2e

package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"

	"github.com/MAVIKE/yad-backend/internal/consts"
	"github.com/MAVIKE/yad-backend/internal/domain"
)

func testGetOrder(s *APITestSuite, jwt string, order *domain.Order) {
	req, err := http.NewRequest("GET", "/api/v1/orders/5", nil)
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

	s.Require().Equal(order.Id, respOrder.Id)
	s.Require().Equal(order.UserId, respOrder.UserId)
	s.Require().Equal(order.RestaurantId, respOrder.RestaurantId)
	s.Require().Equal(order.CourierId, respOrder.CourierId)
	s.Require().Equal(order.DeliveryPrice, respOrder.DeliveryPrice)
	s.Require().Equal(order.TotalPrice, respOrder.TotalPrice)
	s.Require().Equal(order.Status, respOrder.Status)
	if order.Status == consts.OrderCreated {
		s.Require().Equal(order.Paid, respOrder.Paid)
	}
}

func (s *APITestSuite) TestE2EOk() {
	// Create order
	clientId := 1
	clientType := userType
	jwt, err := s.getJWT(clientId, clientType)
	s.NoError(err)

	expectedId := 5
	reqBody := `{"restaurant_id":1}`
	req, err := http.NewRequest("POST", "/api/v1/orders/", bytes.NewBuffer([]byte(reqBody)))
	if err != nil {
		s.FailNow("Failed to build request", err)
	}

	req.Header.Set("Authorization", "Bearer "+jwt)
	req.Header.Set("Content-type", "application/json")

	resp := httptest.NewRecorder()
	s.app.ServeHTTP(resp, req)

	s.Require().Equal(http.StatusOK, resp.Result().StatusCode)

	respData, err := ioutil.ReadAll(resp.Body)
	s.NoError(err)
	var respBody struct {
		Id int `json:"id"`
	}
	err = json.Unmarshal(respData, &respBody)
	s.NoError(err)

	s.Require().Equal(expectedId, respBody.Id)

	// Get order
	order := domain.Order{
		Id:            5,
		UserId:        1,
		RestaurantId:  1,
		CourierId:     0,
		DeliveryPrice: 0,
		TotalPrice:    0,
		Status:        0,
		Paid:          nil,
	}
	testGetOrder(s, jwt, &order)

	// Create first order item
	expectedId = 7
	reqBody = `{"menu_item_id":1, "count":2}`
	req, err = http.NewRequest("POST", "/api/v1/orders/5/items/", bytes.NewBuffer([]byte(reqBody)))
	if err != nil {
		s.FailNow("Failed to build request", err)
	}

	req.Header.Set("Authorization", "Bearer "+jwt)
	req.Header.Set("Content-type", "application/json")

	resp = httptest.NewRecorder()
	s.app.ServeHTTP(resp, req)
	s.Require().Equal(http.StatusOK, resp.Result().StatusCode)

	respData, err = ioutil.ReadAll(resp.Body)
	s.NoError(err)

	err = json.Unmarshal(respData, &respBody)
	s.NoError(err)

	s.Require().Equal(expectedId, respBody.Id)

	// Get first order item
	orderItems := []domain.OrderItem{
		{
			Id:         7,
			OrderId:    5,
			MenuItemId: 1,
			Count:      2,
		},
	}

	req, err = http.NewRequest("GET", "/api/v1/orders/5/items/7", nil)
	if err != nil {
		s.FailNow("Failed to build request", err)
	}
	req.Header.Set("Authorization", "Bearer "+jwt)

	resp = httptest.NewRecorder()
	s.app.ServeHTTP(resp, req)

	s.Require().Equal(http.StatusOK, resp.Result().StatusCode)

	var respOrderItem domain.OrderItem
	respData, err = ioutil.ReadAll(resp.Body)
	s.NoError(err)

	err = json.Unmarshal(respData, &respOrderItem)
	s.NoError(err)

	s.Require().Equal(orderItems[0].Id, respOrderItem.Id)
	s.Require().Equal(orderItems[0].OrderId, respOrderItem.OrderId)
	s.Require().Equal(orderItems[0].MenuItemId, respOrderItem.MenuItemId)
	s.Require().Equal(orderItems[0].Count, respOrderItem.Count)

	// Get order
	order.TotalPrice = 200
	testGetOrder(s, jwt, &order)

	// Update first order item
	reqBody = `{"count":3}`
	updateOrderItemQueryString := fmt.Sprintf("/api/v1/orders/5/items/%d", expectedId)
	req, err = http.NewRequest("PUT", updateOrderItemQueryString, bytes.NewBuffer([]byte(reqBody)))
	if err != nil {
		s.FailNow("Failed to build request", err)
	}

	req.Header.Set("Authorization", "Bearer "+jwt)
	req.Header.Set("Content-type", "application/json")

	resp = httptest.NewRecorder()
	s.app.ServeHTTP(resp, req)

	s.Require().Equal(http.StatusOK, resp.Result().StatusCode)

	// Get order
	order.TotalPrice = 300
	testGetOrder(s, jwt, &order)

	// Create second order item
	expectedId = 8
	reqBody = `{"menu_item_id":2, "count":4}`
	req, err = http.NewRequest("POST", "/api/v1/orders/5/items/", bytes.NewBuffer([]byte(reqBody)))
	if err != nil {
		s.FailNow("Failed to build request", err)
	}

	req.Header.Set("Authorization", "Bearer "+jwt)
	req.Header.Set("Content-type", "application/json")

	resp = httptest.NewRecorder()
	s.app.ServeHTTP(resp, req)

	s.Require().Equal(http.StatusOK, resp.Result().StatusCode)

	respData, err = ioutil.ReadAll(resp.Body)
	s.NoError(err)

	err = json.Unmarshal(respData, &respBody)
	s.NoError(err)

	s.Require().Equal(expectedId, respBody.Id)

	// Get all order items
	orderItems[0].Count = 3
	var secondOrderItem = domain.OrderItem{
		Id:         8,
		OrderId:    5,
		MenuItemId: 2,
		Count:      4,
	}
	orderItems = append(orderItems, secondOrderItem)

	req, err = http.NewRequest("GET", "/api/v1/orders/5/items/", nil)
	if err != nil {
		s.FailNow("Failed to build request", err)
	}
	req.Header.Set("Authorization", "Bearer "+jwt)

	resp = httptest.NewRecorder()
	s.app.ServeHTTP(resp, req)

	s.Require().Equal(http.StatusOK, resp.Result().StatusCode)

	respData, err = ioutil.ReadAll(resp.Body)
	s.NoError(err)

	var respOrderItems = []domain.OrderItem{}
	err = json.Unmarshal(respData, &respOrderItems)
	s.NoError(err)

	for i := 0; i < len(orderItems); i++ {
		s.Require().Equal(orderItems[i].Id, respOrderItems[i].Id)
		s.Require().Equal(orderItems[i].OrderId, respOrderItems[i].OrderId)
		s.Require().Equal(orderItems[i].MenuItemId, respOrderItems[i].MenuItemId)
		s.Require().Equal(orderItems[i].Count, respOrderItems[i].Count)
	}

	// Get order
	order.TotalPrice = 1100
	testGetOrder(s, jwt, &order)

	// Delete Order Item
	req, err = http.NewRequest("DELETE", "/api/v1/orders/5/items/7", nil)
	if err != nil {
		s.FailNow("Failed to build request", err)
	}
	req.Header.Set("Authorization", "Bearer "+jwt)

	resp = httptest.NewRecorder()
	s.app.ServeHTTP(resp, req)

	s.Require().Equal(http.StatusOK, resp.Result().StatusCode)

	// Get order
	order.TotalPrice = 800
	testGetOrder(s, jwt, &order)
}
