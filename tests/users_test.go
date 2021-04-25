package tests

import (
	"bytes"
	"encoding/json"
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
