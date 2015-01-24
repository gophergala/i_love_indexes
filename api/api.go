package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/codegangsta/negroni"
	"github.com/gorilla/mux"
)

func NewAPI() http.Handler {
	router := mux.NewRouter()
	router.HandleFunc("/api/indices", AddIndexOf).Methods("POST")
	stack := negroni.New(negroni.NewLogger(), negroni.NewRecovery())
	stack.UseHandler(router)
	return stack
}

type UnprocessableEntityError struct {
	errors map[string][]string `json:"errors"`
}

func NewUnprocessableEntityError() *UnprocessableEntityError {
	err := &UnprocessableEntityError{}
	err.errors = make(map[string][]string)
	return err
}

func (err *UnprocessableEntityError) Add(field string, errorMessage string) {
	err.errors[field] = append(err.errors[field], errorMessage)
}

type AddIndexOfParams struct {
	URL string `json:"url"`
}

func AddIndexOf(res http.ResponseWriter, req *http.Request) {
	var params *AddIndexOfParams
	err := json.NewDecoder(req.Body).Decode(&params)
	if err != nil {
		res.WriteHeader(500)
		fmt.Fprint(res, err)
		return
	}

	if params.URL == "" {
		res.WriteHeader(422)
		err := NewUnprocessableEntityError()
		err.Add("url", "is empty")
	}
}
