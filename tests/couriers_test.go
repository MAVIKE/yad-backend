package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/MAVIKE/yad-backend/internal/domain"
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

func (s *APITestSuite) TestCourierSignUpError_NoUniquePhone() {
	userId := 1
	clientType := userType

	jwt, err := s.getJWT(userId, clientType)
	s.NoError(err)

	name, phone, password, email, address, working_status := "test_courier", "71234567891", "test_password", "test_courier@test.ru", "{\"latitude\":45,\"longitude\":42}", 0
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

func (s *APITestSuite) TestCourierSignUpError_EmptyRequiredFields() {
	userId := 1
	clientType := userType

	jwt, err := s.getJWT(userId, clientType)
	s.NoError(err)

	name, phone, password, email, address, working_status := "test_courier", "71234567899", "", "test_courier@test.ru", "{\"latitude\":45,\"longitude\":42}", 0
	reqBody := fmt.Sprintf(`{"name":"%s","phone":"%s","password":"%s","email":"%s","address":%s,"working_status":%d}`, name, phone, password, email, address, working_status)

	req, err := http.NewRequest("POST", "/api/v1/couriers/sign-up", bytes.NewBuffer([]byte(reqBody)))
	if err != nil {
		s.FailNow("Failed to build request", err)
	}

	req.Header.Set("Content-type", "application/json")
	req.Header.Set("Authorization", "Bearer "+jwt)

	resp := httptest.NewRecorder()
	s.app.ServeHTTP(resp, req)

	s.Require().Equal(http.StatusBadRequest, resp.Result().StatusCode)
}

func (s *APITestSuite) TestCourierSignInOk() {
	phone, password := "71234567891", "password"
	reqBody := fmt.Sprintf(`{"phone":"%s","password":"%s"}`, phone, password)
	req, err := http.NewRequest("POST", "/api/v1/couriers/sign-in", bytes.NewBuffer([]byte(reqBody)))
	if err != nil {
		s.FailNow("Failed to build request", err)
	}
	req.Header.Set("Content-type", "application/json")

	resp := httptest.NewRecorder()
	s.app.ServeHTTP(resp, req)

	s.Require().Equal(http.StatusOK, resp.Result().StatusCode)
}

func (s *APITestSuite) TestCourierSignInError_WrongPassword() {
	phone, password := "71234567892", "wrong_password"
	reqBody := fmt.Sprintf(`{"phone":"%s","password":"%s"`, phone, password)
	req, err := http.NewRequest("POST", "/api/v1/couriers/sign-in", bytes.NewBuffer([]byte(reqBody)))
	if err != nil {
		s.FailNow("Failed to build request", err)
	}
	req.Header.Set("Content-type", "application/json")

	resp := httptest.NewRecorder()
	s.app.ServeHTTP(resp, req)

	s.Require().Equal(http.StatusBadRequest, resp.Result().StatusCode)
}

func (s *APITestSuite) TestCourierSignInError_NotExists() {
	phone, password := "71234567899", "password"
	reqBody := fmt.Sprintf(`{"phone":"%s","password":"%s"`, phone, password)
	req, err := http.NewRequest("POST", "/api/v1/couriers/sign-in", bytes.NewBuffer([]byte(reqBody)))
	if err != nil {
		s.FailNow("Failed to build request", err)
	}
	req.Header.Set("Content-type", "application/json")

	resp := httptest.NewRecorder()
	s.app.ServeHTTP(resp, req)

	s.Require().Equal(http.StatusBadRequest, resp.Result().StatusCode)
}

func (s *APITestSuite) TestCourierGetOk() {
	userId := 1
	clientType := courierType

	jwt, err := s.getJWT(userId, clientType)
	s.NoError(err)

	req, err := http.NewRequest("GET", "/api/v1/couriers/1", nil)
	if err != nil {
		s.FailNow("Failed to build request", err)
	}

	req.Header.Set("Authorization", "Bearer "+jwt)
	resp := httptest.NewRecorder()
	s.app.ServeHTTP(resp, req)

	s.Require().Equal(http.StatusOK, resp.Result().StatusCode)

	var user domain.Courier
	respData, err := ioutil.ReadAll(resp.Body)
	s.NoError(err)

	err = json.Unmarshal(respData, &user)
	s.NoError(err)

	s.Require().Equal(couriers[0].Id, user.Id)
	s.Require().Equal(couriers[0].Name, user.Name)
	s.Require().Equal(couriers[0].Phone, user.Phone)
	s.Require().Equal("", user.Password)
	s.Require().Equal(couriers[0].Email, user.Email)
	s.Require().Equal(couriers[0].Address, user.Address)
	s.Require().Equal(couriers[0].WorkingStatus, user.WorkingStatus)
}

func (s *APITestSuite) TestCourierGetError_WrongId() {
	userId := 99
	clientType := courierType

	jwt, err := s.getJWT(userId, clientType)
	s.NoError(err)

	req, err := http.NewRequest("GET", "/api/v1/couriers/99", nil)
	if err != nil {
		s.FailNow("Failed to build request", err)
	}

	req.Header.Set("Authorization", "Bearer "+jwt)
	resp := httptest.NewRecorder()
	s.app.ServeHTTP(resp, req)

	s.Require().Equal(http.StatusInternalServerError, resp.Result().StatusCode)
}
