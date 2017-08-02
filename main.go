package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func main() {
	port := 8080

	http.HandleFunc("/helloworld", helloWorldHandler)

	log.Printf("Server starting on port %v\n", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%v", port), nil))
}

func helloWorldHandler(w http.ResponseWriter, r *http.Request) {
	response := helloWorldResponse{Message: "Hello world!"}
	data, err := json.Marshal(response)
	if err != nil {
		panic("Ooops")
	}

	fmt.Fprint(w, string(data))
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
