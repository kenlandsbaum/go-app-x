package rpcclient

import (
	"go-app-x/internal/pb/user"
	"strings"
	"testing"
)

func TestMarshalUnmarshalUser(t *testing.T) {
	expectedUser := &user.User{FirstName: "ken", LastName: "lee", Email: "k@l.com", Id: 1}

	bts, err := MarshalUser(expectedUser)
	if err != nil {
		t.Errorf("got %s but expected nil error\n", err)
	}
	if !strings.Contains(string(bts), "\"first_name\":\"ken\"") {
		t.Errorf("expected first_name but got %s\n", string(bts))
	}

	actualUser, err := UnmarshalUser(bts)
	if err != nil {
		t.Errorf("got %s but expected nil error\n", err)
	}
	if actualUser.FirstName != expectedUser.FirstName {
		t.Errorf("got %s but expected %s\n", actualUser.FirstName, expectedUser.FirstName)
	}
	if actualUser.LastName != expectedUser.LastName {
		t.Errorf("got %s but expected %s\n", actualUser.LastName, expectedUser.LastName)
	}
	if actualUser.Email != expectedUser.Email {
		t.Errorf("got %s but expected %s\n", actualUser.Email, expectedUser.Email)
	}
}
