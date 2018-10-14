package server

import (
	"fmt"
	"github.com/mpfilbin/go-password-hasher/src/password"
	"github.com/mpfilbin/go-password-hasher/src/persistence"
	"log"
	"net/http"
	"time"
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

		idChannel := make(chan int)

		go func() {
			dataStore := persistence.GetInstance()
			id := dataStore.Insert("Placeholder")
			idChannel <- id

			time.Sleep(5 * time.Second)

			encoded := password.Encode(request.FormValue("password"))
			dataStore.Update(id, encoded)
			log.Println("Persistence of encoded password complete")

		}()

		id := <- idChannel
		message := fmt.Sprintf("ID will be: %v", id)
		response.WriteHeader(http.StatusAccepted)
		response.Write([]byte(message))
	default:
		http.NotFound(response, request)
	}
}

func Listen(port int) {
	serverMux := http.NewServeMux() // Create isolated server mutex

	serverMux.HandleFunc("/hash", encodeAndPersistPassword)

	address :=fmt.Sprintf(":%d", port)
	log.Printf("Listening at %s", address)

	err := http.ListenAndServe(address, serverMux) // set listen port
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}