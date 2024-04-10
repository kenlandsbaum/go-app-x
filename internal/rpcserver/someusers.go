package rpcserver

import "go-app-x/internal/user"

var (
	someUsers = []*user.User{
		{Id: 1, FirstName: "ken", LastName: "lee", Email: "ken@site.com"},
		{Id: 2, FirstName: "jen", LastName: "lee", Email: "jen@site.com"},
		{Id: 3, FirstName: "sam", LastName: "jones", Email: "sam@site.com"},
		{Id: 4, FirstName: "jill", LastName: "smith", Email: "jill@site.com"},
	}
)
