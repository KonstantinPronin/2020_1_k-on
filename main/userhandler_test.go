package main

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

type StubContext struct {
	context.Context
}

func (mc *StubContext) Value(interface{}) interface{} {
	return map[string]string{"id": "10"}
}

func TestUserHandler_Login_EmptyLogin(t *testing.T) {
	userHandler := createUserHandler()
	user := User{
		Password: "test",
	}
	body, err := json.Marshal(user)
	if err != nil {
		t.Fatal(err)
	}

	response := httptest.NewRecorder()
	request, err := http.NewRequest("POST", "/login", bytes.NewBuffer(body))
	if err != nil {
		t.Fatal(err)
	}

	handler := http.HandlerFunc(userHandler.Login)
	handler.ServeHTTP(response, request)
	data, err := ioutil.ReadAll(response.Body)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, "bad login\n", string(data))
	assert.Equal(t, http.StatusBadRequest, response.Code)
}

func TestUserHandler_Login_UserDoesNotExist(t *testing.T) {
	userHandler := createUserHandler()
	user := User{
		Username: "test",
		Password: "test",
	}
	body, err := json.Marshal(user)
	if err != nil {
		t.Fatal(err)
	}

	response := httptest.NewRecorder()
	request, err := http.NewRequest("POST", "/login", bytes.NewBuffer(body))
	if err != nil {
		t.Fatal(err)
	}

	handler := http.HandlerFunc(userHandler.Login)
	handler.ServeHTTP(response, request)
	data, err := ioutil.ReadAll(response.Body)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, "no user\n", string(data))
	assert.Equal(t, http.StatusNotFound, response.Code)
}

func TestUserHandler_Login_WrongPassword(t *testing.T) {
	userHandler := createUserHandler()
	user := User{
		Username: "test",
		Password: "test",
	}
	wrongUser := User{
		Username: "test",
		Password: "wrong",
	}
	userHandler.users.users[user.Username] = &user

	body, err := json.Marshal(wrongUser)
	if err != nil {
		t.Fatal(err)
	}

	response := httptest.NewRecorder()
	request, err := http.NewRequest("POST", "/login", bytes.NewBuffer(body))
	if err != nil {
		t.Fatal(err)
	}

	handler := http.HandlerFunc(userHandler.Login)
	handler.ServeHTTP(response, request)
	data, err := ioutil.ReadAll(response.Body)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, "bad password\n", string(data))
	assert.Equal(t, http.StatusBadRequest, response.Code)
}

func TestUserHandler_Login(t *testing.T) {
	userHandler := createUserHandler()
	user := User{
		Username: "test",
		Password: "test",
	}
	userHandler.users.users[user.Username] = &user
	body, err := json.Marshal(user)
	if err != nil {
		t.Fatal(err)
	}

	response := httptest.NewRecorder()
	request, err := http.NewRequest("POST", "/login", bytes.NewBuffer(body))
	if err != nil {
		t.Fatal(err)
	}

	handler := http.HandlerFunc(userHandler.Login)
	handler.ServeHTTP(response, request)
	_, err = ioutil.ReadAll(response.Body)
	if err != nil {
		t.Fatal(err)
	}

	assert.NotEmpty(t, response.Result().Cookies())
	assert.Equal(t, http.StatusOK, response.Code)
}

func TestUserHandler_Logout_WithoutCookie(t *testing.T) {
	userHandler := createUserHandler()

	response := httptest.NewRecorder()
	request, err := http.NewRequest("POST", "/logout", nil)
	if err != nil {
		t.Fatal(err)
	}

	handler := http.HandlerFunc(userHandler.Logout)
	handler.ServeHTTP(response, request)
	data, err := ioutil.ReadAll(response.Body)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, "no session\n", string(data))
	assert.Equal(t, http.StatusUnauthorized, response.Code)
}

func TestUserHandler_Logout(t *testing.T) {
	userHandler := createUserHandler()
	userHandler.sessions["test"] = 0

	response := httptest.NewRecorder()
	request, err := http.NewRequest("POST", "/logout", nil)
	if err != nil {
		t.Fatal(err)
	}
	cookie := &http.Cookie{
		Name:    "session_id",
		Value:   "test",
		Expires: time.Now().Add(10 * time.Hour),
	}
	request.AddCookie(cookie)

	handler := http.HandlerFunc(userHandler.Logout)
	handler.ServeHTTP(response, request)
	_, err = ioutil.ReadAll(response.Body)
	if err != nil {
		t.Fatal(err)
	}

	cookie = response.Result().Cookies()[0]
	assert.GreaterOrEqual(t, time.Now().AddDate(0, 0, -1).String(), cookie.Expires.String())
	assert.Equal(t, http.StatusOK, response.Code)
}

func TestUserHandler_Add_EmptyLogin(t *testing.T) {
	userHandler := createUserHandler()
	user := User{
		Password: "test",
	}
	body, err := json.Marshal(user)
	if err != nil {
		t.Fatal(err)
	}

	response := httptest.NewRecorder()
	request, err := http.NewRequest("POST", "/signup", bytes.NewBuffer(body))
	if err != nil {
		t.Fatal(err)
	}

	handler := http.HandlerFunc(userHandler.Add)
	handler.ServeHTTP(response, request)
	data, err := ioutil.ReadAll(response.Body)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, "bad parameters\n", string(data))
	assert.Equal(t, http.StatusBadRequest, response.Code)
}

func TestUserHandler_Add(t *testing.T) {
	userHandler := createUserHandler()
	user := User{
		Username: "test",
		Password: "test",
	}
	body, err := json.Marshal(user)
	if err != nil {
		t.Fatal(err)
	}

	response := httptest.NewRecorder()
	request, err := http.NewRequest("POST", "/signup", bytes.NewBuffer(body))
	if err != nil {
		t.Fatal(err)
	}

	handler := http.HandlerFunc(userHandler.Add)
	handler.ServeHTTP(response, request)
	_, err = ioutil.ReadAll(response.Body)
	if err != nil {
		t.Fatal(err)
	}

	_, ok := userHandler.users.GetByName("test")
	assert.True(t, ok)
	assert.Equal(t, http.StatusOK, response.Code)
}

func TestUserHandler_UserById_WrongId(t *testing.T) {
	userHandler := createUserHandler()
	response := httptest.NewRecorder()
	request, err := http.NewRequest("GET", "/user/qwe", nil)
	if err != nil {
		t.Fatal(err)
	}

	handler := http.HandlerFunc(userHandler.UserById)
	handler.ServeHTTP(response, request)
	data, err := ioutil.ReadAll(response.Body)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, "bad id\n", string(data))
	assert.Equal(t, http.StatusBadRequest, response.Code)
}

func TestUserHandler_UserById_UserDoesNotExist(t *testing.T) {
	userHandler := createUserHandler()
	response := httptest.NewRecorder()
	request, err := http.NewRequest("GET", "/user", nil)
	if err != nil {
		t.Fatal(err)
	}
	ctx := new(StubContext)
	request = request.WithContext(ctx)

	handler := http.HandlerFunc(userHandler.UserById)
	handler.ServeHTTP(response, request)
	data, err := ioutil.ReadAll(response.Body)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, "no user\n", string(data))
	assert.Equal(t, 404, response.Code)
}

func TestUserHandler_UserById(t *testing.T) {
	userHandler := createUserHandler()
	user := User{
		Id:       10,
		Username: "test",
		Password: "test",
	}
	userHandler.users.users[user.Username] = &user
	response := httptest.NewRecorder()
	request, err := http.NewRequest("GET", "/user", nil)
	if err != nil {
		t.Fatal(err)
	}
	ctx := new(StubContext)
	request = request.WithContext(ctx)

	handler := http.HandlerFunc(userHandler.UserById)
	handler.ServeHTTP(response, request)
	data, err := ioutil.ReadAll(response.Body)
	if err != nil {
		t.Fatal(err)
	}
	result := new(User)
	err = json.Unmarshal(data, &result)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, user, *result)
	assert.Equal(t, http.StatusOK, response.Code)
}
