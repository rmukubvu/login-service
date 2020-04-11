package store

import (
	"crypto/sha1"
	"encoding/hex"
	"errors"
	"github.com/mongodb/mongo-go-driver/bson/primitive"
	"github.com/mongodb/mongo-go-driver/mongo"
	"strings"
)

var (
	cache                 = make(map[string]User)
	userNotExistError     = errors.New("user does not exist")
	passwordMismatchError = errors.New("user/password combination is invalid")
	client                *mongo.Client
	db                    *mongo.Database
)

const (
	schemaName          = "amakhosi_logins"
	loginCollectionName = "logins"
	databaseUri         = "mongodb://localhost:27017"
)

type User struct {
	UserId    primitive.ObjectID `json:"_id" bson:"_id,omitempty"`
	UserName  string             `json:"username" bson:"username,omitempty"`
	Password  string             `json:"password" bson:"password,omitempty"`
	IsEnabled int                `json:"isEnabled" bson:"is_enabled,omitempty"`
	FirstName string             `json:"firstName" bson:"first_name,omitempty"`
	LastName  string             `json:"lastName" bson:"last_name,omitempty"`
}

type Authenticate struct {
	UserName string `json:"username"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Response string `json:"response"`
	IsError  bool   `json:"error"`
	User
}

//do an update if the record exists
func ValidateLogin(userName, passWord string) LoginResponse {
	//fetch the user first
	u, err := FetchRecord(userName)
	//_, exists, u := Search(userName)
	if err != nil {
		return LoginResponse{
			Response: userNotExistError.Error(),
			IsError:  true,
			User:     User{},
		}
	}

	var isError = false
	var pwdError = "ok"
	//hash and check the hashes
	hashedPassword := getHashedPassword(passWord)
	if strings.Compare(u.Password, hashedPassword) != 0 {
		isError = true
		pwdError = passwordMismatchError.Error()
		u = User{}
	}
	return LoginResponse{
		Response: pwdError,
		IsError:  isError,
		User:     u,
	}
}

func Search(userName string) (canInsert, exists bool, u User) {
	//get the user from the map
	u = cache[userName]
	if u == (User{}) {
		//can insert
		canInsert = true
		exists = false
		return
	}
	canInsert = false
	exists = true
	return
}

func (u *User) AddToMap() {
	cache[u.UserName] = *u
}

func loadMap() {
	if users := fetchAllRecords(); users != nil {
		for _, u := range users {
			cache[u.UserName] = u
		}
	}
}

func getHashedPassword(password string) string {
	h := sha1.New()
	h.Write([]byte(password))
	bs := h.Sum(nil)
	return hex.EncodeToString(bs)
}
