package repository

import "github.com/rmukubvu/login-service/store"

type UserEntity struct {
	store.User
}

func (u *UserEntity) Add() error {
	return u.InsertRecord()
}

func (u *UserEntity) Update() error {
	return u.UpdateRecord()
}

func Search(userName string) (store.User, error) {
	return store.FetchRecord(userName)
}

func AuthenticateLogin(u store.Authenticate) store.LoginResponse {
	return validateLogin(u.UserName, u.Password)
}

func validateLogin(userName, passWord string) store.LoginResponse {
	return store.ValidateLogin(userName, passWord)
}
