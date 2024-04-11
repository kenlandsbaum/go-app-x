package main

import (
	"context"
	"encoding/json"
	"errors"
	"go-app-x/internal/router"
	"go-app-x/internal/user"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"google.golang.org/grpc"
)

var (
	mockUser         *user.User
	mockUserResponse *user.NewUserResponse
)

/*
GetAll(ctx context.Context, in *AllUsersRequest, opts ...grpc.CallOption) (UserService_GetAllClient, error)
GetOne(ctx context.Context, in *UserRequest, opts ...grpc.CallOption) (*User, error)
Create(ctx context.Context, in *NewUserRequest, opts ...grpc.CallOption) (*NewUserResponse, error)
*/

type mockRpcClient struct{}

func (m mockRpcClient) GetAll(
	ctx context.Context,
	in *user.AllUsersRequest,
	opts ...grpc.CallOption) (user.UserService_GetAllClient, error) {
	return nil, nil
}

func (m mockRpcClient) GetOne(
	ctx context.Context,
	in *user.UserRequest,
	opts ...grpc.CallOption) (*user.User, error) {
	if mockUser != nil {
		return mockUser, nil
	}
	return nil, errors.New("no user")
}

func (m mockRpcClient) Create(
	ctx context.Context,
	in *user.NewUserRequest,
	opts ...grpc.CallOption) (*user.NewUserResponse, error) {
	if mockUserResponse != nil {
		return mockUserResponse, nil
	}
	return nil, errors.New("no user response")
}

func getTestWebServer() *webserver {
	ctx := context.Background()
	testClient := mockRpcClient{}
	r := router.New()
	return createServer(ctx, testClient, r)
}

func TestWebServer_getUser(t *testing.T) {
	testWebServer := getTestWebServer()

	mockUser = &user.User{Id: 1, FirstName: "ken", LastName: "lee", Email: "k@l.com"}
	w1 := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/user", nil)

	testWebServer.getUser(w1, req)
	result := w1.Result()
	if result.StatusCode != http.StatusOK {
		t.Errorf("got %d but expected %d", result.StatusCode, http.StatusOK)
	}
	defer result.Body.Close()
	bts, _ := io.ReadAll(result.Body)

	var actualUser user.User
	json.Unmarshal(bts, &actualUser)

	if actualUser.FirstName != mockUser.FirstName {
		t.Errorf("got %s but expected %s\n", actualUser.FirstName, mockUser.FirstName)
	}

	mockUser = nil
	w2 := httptest.NewRecorder()
	req, _ = http.NewRequest(http.MethodGet, "/user", nil)

	testWebServer.getUser(w2, req)
	result = w2.Result()
	if result.StatusCode != http.StatusInternalServerError {
		t.Errorf("got %d but expected %d", result.StatusCode, http.StatusInternalServerError)
	}
	defer result.Body.Close()
	bts, _ = io.ReadAll(result.Body)

	if string(bts) != "no user" {
		t.Errorf("got %s but expected %s\n", string(bts), "no user")
	}
}

func TestWebServer_addUser(t *testing.T) {
	testWebServer := getTestWebServer()

	reqBody := strings.NewReader(`{"first_name":"ken","last_name":"test","email":"k@l.com"}`)

	w1 := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/user", reqBody)

	mockUserResponse = &user.NewUserResponse{Id: 1}
	testWebServer.addUser(w1, req)
	result := w1.Result()
	if result.StatusCode != http.StatusCreated {
		t.Errorf("got %d but expected %d", result.StatusCode, http.StatusOK)
	}
	defer result.Body.Close()
	bts, _ := io.ReadAll(result.Body)

	if !strings.Contains(string(bts), "numberOfUsers") {
		t.Errorf("got %s but expected a response with number of current users", string(bts))
	}
}
