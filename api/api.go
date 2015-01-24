package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/GopherGala/i_love_indexes/elasticsearch"
	"github.com/Scalingo/go-workers"
	"github.com/codegangsta/negroni"
	"github.com/gorilla/mux"
	"gopkg.in/errgo.v1"
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
		res.WriteHeader(400)
		fmt.Fprintf(res, "invalid JSON: %v", err)
		return
	}

	if params.URL == "" {
		res.WriteHeader(422)
		err := NewUnprocessableEntityError()
		err.Add("url", "can't be blank")
		json.NewEncoder(res).Encode(&err)
		return
	}

	index := &elasticsearch.IndexOf{URL: params.URL}
	err = elasticsearch.Index(index)
	if err != nil {
		res.WriteHeader(500)
		fmt.Fprintf(res, errgo.Details(err))
		return
	}

	workers.Enqueue("index-crawler", "CrawlWorker", []string{index.URL, index.Id})
	res.WriteHeader(201)
	json.NewEncoder(res).Encode(&index)
}
