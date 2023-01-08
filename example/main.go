package main

import (
	"context"
	"net/http"

	. "github.com/danielvladco/generic-http-handlers"
)

func main() {
	svc := accountService{}

	accountHandler := MakeAccountHandler(svc)

	http.ListenAndServe(":8080", accountHandler)
}

func MakeAccountHandler(svc AccountService) http.Handler {
	mux := http.NewServeMux()
	mux.Handle("/account", GET(svc.GetAccount))
	mux.Handle("/account", POST(svc.CreateAccount))
	mux.Handle("/account", DELETE(svc.DeleteAccount))
	return mux
}

type Account struct {
	ID    string `json:"id"`
	Email string `json:"email"`
}

type AccountByID struct {
	ID string `json:"id"`
}

type AccountService interface {
	CreateAccount(ctx context.Context, account *Account) (*Account, error)
	GetAccount(ctx context.Context, byID *AccountByID) (*Account, error)
	DeleteAccount(ctx context.Context, byID *AccountByID) (bool, error)
}

type accountService struct {
}

func (a accountService) CreateAccount(ctx context.Context, account *Account) (*Account, error) {
	//TODO implement me
	panic("implement me")
}

func (a accountService) GetAccount(ctx context.Context, byID *AccountByID) (*Account, error) {
	//TODO implement me
	panic("implement me")
}

func (a accountService) DeleteAccount(ctx context.Context, byID *AccountByID) (bool, error) {
	//TODO implement me
	panic("implement me")
}
