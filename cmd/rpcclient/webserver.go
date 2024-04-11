package main

import (
	"context"
	"fmt"
	"go-app-x/internal/pb/user"
	"go-app-x/internal/router"
	"go-app-x/internal/rpcclient"
	"io"
	"math/rand/v2"
	"net/http"
	"os"
)

type webserver struct {
	ctx       context.Context
	rpcClient user.UserServiceClient
	router    *router.Router
	temp      int64
}

func createServer(ctx context.Context, rpcClient user.UserServiceClient, router *router.Router) *webserver {
	return &webserver{ctx: ctx, rpcClient: rpcClient, router: router, temp: 4}
}

func (ws *webserver) routes() {
	ws.router.Get("/user", withJsonContent(ws.getUser))
	ws.router.Post("/user", withJsonContent(ws.addUser))
}

func (ws *webserver) getUser(w http.ResponseWriter, r *http.Request) {
	id := rand.Int64N(ws.temp) + 1
	res, err := ws.rpcClient.GetOne(ws.ctx, &user.UserRequest{Id: id})
	if err != nil {
		handleFacePlant(w, err)
		return
	}
	userBts, err := rpcclient.MarshalUser(res)
	if err != nil {
		handleFacePlant(w, err)
		return
	}
	handleSuccess(w, userBts, http.StatusOK)
}

func (ws *webserver) addUser(w http.ResponseWriter, r *http.Request) {
	body := r.Body
	bts, err := io.ReadAll(body)
	if err != nil {
		handleFacePlant(w, err)
		return
	}
	newUser, err := rpcclient.UnmarshalUser(bts)
	if err != nil {
		handleFacePlant(w, err)
		return
	}
	res, err := ws.rpcClient.Create(
		ws.ctx,
		&user.NewUserRequest{
			FirstName: newUser.FirstName,
			LastName:  newUser.LastName,
			Email:     newUser.Email,
		})
	if err != nil {
		handleFacePlant(w, err)
		return
	}
	ws.temp += 1
	resBts := []byte(fmt.Sprintf(`{"numberOfUsers": %d, "newUserId": %d}`, ws.temp, res.Id))
	handleSuccess(w, resBts, http.StatusCreated)
}

func (ws *webserver) start() {
	if err := ws.router.Start(os.Getenv("WEB_ADDRESS")); err != nil {
		panic(err)
	}
}
