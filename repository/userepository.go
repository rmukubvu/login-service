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

func Search(userName string) (store.User,error){
	return store.FetchRecord(userName)
}

func ValidateLogin(userName,passWord string) store.LoginResponse {
	return store.ValidateLogin(userName,passWord)
}