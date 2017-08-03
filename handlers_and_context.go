package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func main() {
	port := 8080

	var helloWorldHandler = newValidationHandler(NewHelloWorldHandler())
	http.Handle("/helloworld", helloWorldHandler)
	log.Printf("Server starting on port %v\n", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%v", port), nil))
}

type helloWorldResponse struct {
	//change the output field name to message
	Message string `json:"message"`
	//Do not output the field
	Author string `json:"-"`
	//Do not output the field if it is empty
	Date string `json:",omitempty"`
	//Convert output to a string
	Id int `json:"id,string"`
}

type helloWorldRequest struct {
	Name string `json:"name"`
}

type validationHandler struct {
	next http.Handler
}

func newValidationHandler(next http.Handler) http.Handler {
	return validationHandler{next: next}
}

type validationContextKey string

func (h validationHandler) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	var request helloWorldRequest
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&request)
	if err != nil {
		http.Error(rw, "Bad Request", http.StatusBadRequest)
		return
	}
	c := context.WithValue(r.Context(), validationContextKey("name"), request.Name)
	r = r.WithContext(c)
	h.next.ServeHTTP(rw, r)
}

type helloWorldHandler struct{}

func NewHelloWorldHandler() http.Handler {
	return helloWorldHandler{}
}

func (h helloWorldHandler) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	name := r.Context().Value(validationContextKey("name")).(string)
	response := helloWorldResponse{Message: "Hello " + name}
	encoder := json.NewEncoder(rw)
	encoder.Encode(&response)
}
