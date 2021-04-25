package tests

import (
	"bytes"
	"encoding/json"
	"github.com/MAVIKE/yad-backend/internal/domain"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
)

func (s *APITestSuite) TestUserSignUpOk() {
	expectedId := 4
	reqBody := `{"name":"test_name","phone":"78880001119","password":"test_password","email":"test@test.ru","address":{"latitude":45,"longitude":42}}`
	req, err := http.NewRequest("POST", "/api/v1/users/sign-up", bytes.NewBuffer([]byte(reqBody)))
	if err != nil {
		s.FailNow("Failed to build request", err)
	}
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
}

func (s *APITestSuite) TestUserSignUpError_NotUniquePhone() {
	reqBody := `{"name":"test_name","phone":"71234567890","password":"test_password","email":"test@test.ru","address":{"latitude":45,"longitude":42}}`
	req, err := http.NewRequest("POST", "/api/v1/users/sign-up", bytes.NewBuffer([]byte(reqBody)))
	if err != nil {
		s.FailNow("Failed to build request", err)
	}
	req.Header.Set("Content-type", "application/json")

	resp := httptest.NewRecorder()
	s.app.ServeHTTP(resp, req)

	s.Require().Equal(http.StatusInternalServerError, resp.Result().StatusCode)
}

func (s *APITestSuite) TestUserSignUpError_EmptyRequiredFields() {
	reqBody := `{"name":"test_name","phone":"","password":"","email":"test@test.ru","address":{"latitude":45,"longitude":42}}`
	req, err := http.NewRequest("POST", "/api/v1/users/sign-up", bytes.NewBuffer([]byte(reqBody)))
	if err != nil {
		s.FailNow("Failed to build request", err)
	}
	req.Header.Set("Content-type", "application/json")

	resp := httptest.NewRecorder()
	s.app.ServeHTTP(resp, req)

	s.Require().Equal(http.StatusBadRequest, resp.Result().StatusCode)
}

func (s *APITestSuite) TestUserSignInOk() {
	reqBody := `{"phone":"71234567890","password":"password"}`
	req, err := http.NewRequest("POST", "/api/v1/users/sign-in", bytes.NewBuffer([]byte(reqBody)))
	if err != nil {
		s.FailNow("Failed to build request", err)
	}
	req.Header.Set("Content-type", "application/json")

	resp := httptest.NewRecorder()
	s.app.ServeHTTP(resp, req)

	s.Require().Equal(http.StatusOK, resp.Result().StatusCode)
}

func (s *APITestSuite) TestUserSignInError_WrongPassword() {
	reqBody := `{"phone":"71234567890","password":"wrong_password"}`
	req, err := http.NewRequest("POST", "/api/v1/users/sign-in", bytes.NewBuffer([]byte(reqBody)))
	if err != nil {
		s.FailNow("Failed to build request", err)
	}
	req.Header.Set("Content-type", "application/json")

	resp := httptest.NewRecorder()
	s.app.ServeHTTP(resp, req)

	s.Require().Equal(http.StatusInternalServerError, resp.Result().StatusCode)
}

func (s *APITestSuite) TestUserSignInError_NotExists() {
	reqBody := `{"phone":"71234567899","password":"password"}`
	req, err := http.NewRequest("POST", "/api/v1/users/sign-in", bytes.NewBuffer([]byte(reqBody)))
	if err != nil {
		s.FailNow("Failed to build request", err)
	}
	req.Header.Set("Content-type", "application/json")

	resp := httptest.NewRecorder()
	s.app.ServeHTTP(resp, req)

	s.Require().Equal(http.StatusInternalServerError, resp.Result().StatusCode)
}

func (s *APITestSuite) TestUserGetOk() {
	userId := 1
	clientType := userType

	jwt, err := s.getJWT(userId, clientType)
	s.NoError(err)

	req, err := http.NewRequest("GET", "/api/v1/users/1", nil)
	if err != nil {
		s.FailNow("Failed to build request", err)
	}

	req.Header.Set("Authorization", "Bearer "+jwt)
	resp := httptest.NewRecorder()
	s.app.ServeHTTP(resp, req)

	s.Require().Equal(http.StatusOK, resp.Result().StatusCode)

	var user domain.User
	respData, err := ioutil.ReadAll(resp.Body)
	s.NoError(err)

	err = json.Unmarshal(respData, &user)
	s.NoError(err)

	s.Require().Equal(users[0].Id, user.Id)
	s.Require().Equal(users[0].Name, user.Name)
	s.Require().Equal(users[0].Phone, user.Phone)
	s.Require().Equal("", user.Password)
	s.Require().Equal(users[0].Email, user.Email)
	s.Require().Equal(users[0].Address, user.Address)
}

func (s *APITestSuite) TestUserGetError_WrongId() {
	userId := 2
	clientType := userType

	jwt, err := s.getJWT(userId, clientType)
	s.NoError(err)

	req, err := http.NewRequest("GET", "/api/v1/users/1", nil)
	if err != nil {
		s.FailNow("Failed to build request", err)
	}

	req.Header.Set("Authorization", "Bearer "+jwt)
	resp := httptest.NewRecorder()
	s.app.ServeHTTP(resp, req)

	s.Require().Equal(http.StatusInternalServerError, resp.Result().StatusCode)
}

func (s *APITestSuite) TestUserGetError_WrongClientType() {
	userId := 1
	clientType := courierType

	jwt, err := s.getJWT(userId, clientType)
	s.NoError(err)

	req, err := http.NewRequest("GET", "/api/v1/users/1", nil)
	if err != nil {
		s.FailNow("Failed to build request", err)
	}

	req.Header.Set("Authorization", "Bearer "+jwt)
	resp := httptest.NewRecorder()
	s.app.ServeHTTP(resp, req)

	s.Require().Equal(http.StatusInternalServerError, resp.Result().StatusCode)
}
