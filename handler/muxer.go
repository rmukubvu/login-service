package handler

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/rmukubvu/amakhosi/model"
	"github.com/rmukubvu/login-service/repository"
	"github.com/rmukubvu/login-service/store"
	"io/ioutil"
	"net/http"
)

func InitRouter() *mux.Router {
	r := mux.NewRouter()
	api := r.PathPrefix("/api/v1").Subrouter()
	api.HandleFunc("/user", add).Methods(http.MethodPost)
	api.HandleFunc("/user", update).Methods(http.MethodPut)
	api.HandleFunc("/user/{id}", search).Methods(http.MethodGet)
	api.HandleFunc("/user/authenticate", authenticate).Methods(http.MethodPost)
	return r
}

func add(response http.ResponseWriter, request *http.Request) {
	response.Header().Add("content-type", "application/json")
	user := repository.UserEntity{}
	reqBody, err := ioutil.ReadAll(request.Body)
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(response).Encode(generateErrorMessage(err.Error()))
		return
	}
	//check if its a valid json string
	if ok := validJson(string(reqBody)); !ok {
		response.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(response).Encode(generateErrorMessage("invalid json string"))
		return
	}
	//continue to unmarshal
	if err = json.Unmarshal(reqBody, &user); err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(response).Encode(generateErrorMessage(err.Error()))
		return
	}
	//add to database
	err = user.InsertRecord()
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(response).Encode(generateErrorMessage(err.Error()))
		return
	}
	response.WriteHeader(http.StatusCreated)
	json.NewEncoder(response).Encode(user)
}

func update(response http.ResponseWriter, request *http.Request) {
	response.Header().Add("content-type", "application/json")
	user := repository.UserEntity{}
	reqBody, err := ioutil.ReadAll(request.Body)
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(response).Encode(generateErrorMessage(err.Error()))
		return
	}
	//check if its a valid json string
	if ok := validJson(string(reqBody)); !ok {
		response.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(response).Encode(generateErrorMessage("invalid json string"))
		return
	}
	//continue to unmarshal
	json.Unmarshal(reqBody, &user)
	//add to database
	err = user.UpdateRecord()
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(response).Encode(generateErrorMessage(err.Error()))
		return
	}
	response.WriteHeader(http.StatusOK)
	json.NewEncoder(response).Encode(user)
}

func search(response http.ResponseWriter, req *http.Request) {
	pathParams := mux.Vars(req)
	response.Header().Set("Content-Type", "application/json")

	if val, ok := pathParams["id"]; !ok {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(generateErrorMessage("id not specified or wrong id format")))
		return
	} else {
		user, err := repository.Search(val)
		if err != nil {
			response.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(response).Encode(generateErrorMessage(err.Error()))
			return
		}
		response.WriteHeader(http.StatusOK)
		json.NewEncoder(response).Encode(user)
	}
}

func authenticate(response http.ResponseWriter, req *http.Request) {
	response.Header().Set("Content-Type", "application/json")
	user := store.Authenticate{}
	reqBody, err := ioutil.ReadAll(req.Body)
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(response).Encode(generateErrorMessage(err.Error()))
		return
	}
	//check if its a valid json string
	if ok := validJson(string(reqBody)); !ok {
		response.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(response).Encode(generateErrorMessage("invalid json string"))
		return
	}
	//continue to unmarshal
	err = json.Unmarshal(reqBody, &user)
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(response).Encode(err.Error())
		return
	}
	result := repository.AuthenticateLogin(user)
	response.WriteHeader(http.StatusOK)
	json.NewEncoder(response).Encode(result)
}

func validJson(str string) bool {
	var js json.RawMessage
	return json.Unmarshal([]byte(str), &js) == nil
}

func generateErrorMessage(e string) string {
	ie := model.InternalError{Message: e}
	buf, err := json.Marshal(ie)
	if err != nil {
		return string([]byte(fmt.Sprintf(`{"message": "%s"}`, e)))
	}
	return string(buf)
}
