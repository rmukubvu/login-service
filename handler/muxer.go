package handler

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/rmukubvu/amakhosi/model"
	"github.com/rmukubvu/login-service/repository"
	"io/ioutil"
	"net/http"
)

func InitRouter() *mux.Router{
	r := mux.NewRouter()
	api := r.PathPrefix("/api/v1").Subrouter()
	api.HandleFunc("/user", add).Methods(http.MethodPost)
	api.HandleFunc("/user", update).Methods(http.MethodPut)
	api.HandleFunc("/user/{id}", search).Methods(http.MethodGet)
	api.HandleFunc("/user/validate/{id}/{pwd}", validate).Methods(http.MethodGet)
	return r
}

func add(response http.ResponseWriter,request *http.Request){
	response.Header().Add("content-type","application/json")
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
	if err = json.Unmarshal(reqBody, &user) ; err != nil {
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

func update(response http.ResponseWriter,request *http.Request){
	response.Header().Add("content-type","application/json")
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
	//continue to unmarshall
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

func search(w http.ResponseWriter, req *http.Request) {
	pathParams := mux.Vars(req)
	w.Header().Set("Content-Type", "application/json")

	if val, ok := pathParams["id"]; !ok {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(generateErrorMessage("id not specified or wrong id format")))
		return
	} else {
		user, err := repository.Search(val)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(generateErrorMessage(err.Error()))
			return
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(user)
	}
}

func validate(w http.ResponseWriter, req *http.Request) {
	pathParams := mux.Vars(req)
	w.Header().Set("Content-Type", "application/json")

	userName:= pathParams["id"]
	password:= pathParams["pwd"]

	response := repository.ValidateLogin(userName,password)
	if response.IsError {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
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