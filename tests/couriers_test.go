package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
)

func (s *APITestSuite) TestCourierSignUpOk() {
	expectedId := 6

	userId := 1
	clientType := adminType

	jwt, err := s.getJWT(userId, clientType)
	s.NoError(err)

	name, phone, password, email, address, working_status := "test_courier", "71234567899", "test_password", "test_courier@test.ru", "{\"latitude\":45,\"longitude\":42}", 0
	reqBody := fmt.Sprintf(`{"name":"%s","phone":"%s","password":"%s","email":"%s","address":%s,"working_status":%d}`, name, phone, password, email, address, working_status)

	req, err := http.NewRequest("POST", "/api/v1/couriers/sign-up", bytes.NewBuffer([]byte(reqBody)))
	if err != nil {
		s.FailNow("Failed to build request", err)
	}

	req.Header.Set("Content-type", "application/json")
	req.Header.Set("Authorization", "Bearer "+jwt)

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
}

func (s *APITestSuite) TestCourierSignUpError_NotAdmin() {
	userId := 1
	clientType := userType

	jwt, err := s.getJWT(userId, clientType)
	s.NoError(err)

	name, phone, password, email, address, working_status := "test_courier", "71234567899", "test_password", "test_courier@test.ru", "{\"latitude\":45,\"longitude\":42}", 0
	reqBody := fmt.Sprintf(`{"name":"%s","phone":"%s","password":"%s","email":"%s","address":%s,"working_status":%d}`, name, phone, password, email, address, working_status)

	req, err := http.NewRequest("POST", "/api/v1/couriers/sign-up", bytes.NewBuffer([]byte(reqBody)))
	if err != nil {
		s.FailNow("Failed to build request", err)
	}

	req.Header.Set("Content-type", "application/json")
	req.Header.Set("Authorization", "Bearer "+jwt)

	resp := httptest.NewRecorder()
	s.app.ServeHTTP(resp, req)

	s.Require().Equal(http.StatusInternalServerError, resp.Result().StatusCode)
}
