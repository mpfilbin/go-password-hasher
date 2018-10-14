package server

import (
	"fmt"
	"github.com/mpfilbin/go-password-hasher/src/password"
	"log"
	"net/http"
)

func encodeAndPersistPassword(response http.ResponseWriter, request *http.Request) {
	log.Printf("Received %s /hash", request.Method)
	switch request.Method {
	case "POST":
		if err := request.ParseForm(); err != nil {
			log.Printf("Error - Unable to parse form data %v", err)
			http.Error(response, "Invalid form data", http.StatusBadRequest)
			return
		}

		fmt.Fprintf(response, "Post from website! r.PostFrom = %v\n", request.PostForm)
		clearTextPassword := request.FormValue("password")
		log.Printf("Received request for %s", clearTextPassword)

		response.Write([]byte(password.Encode(clearTextPassword)))
	default:
		http.NotFound(response, request)
	}
}

func Listen(port int) {
	http.HandleFunc("/hash", encodeAndPersistPassword)

	address :=fmt.Sprintf(":%d", port)
	log.Printf("Listening at %s", address)

	err := http.ListenAndServe(address, nil) // set listen port
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}