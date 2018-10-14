package server

import (
	"encoding/json"
	"fmt"
	"github.com/mpfilbin/go-password-hasher/src/password"
	"github.com/mpfilbin/go-password-hasher/src/persistence"
	"log"
	"net/http"
	"time"
)

type persistenceResult struct {
	ID            int `json:"id"`
	TimeAvailable string `json:"timeAvailable"`
	URL           string `json:"url"`

}

func encodeAndPersistPassword(response http.ResponseWriter, request *http.Request) {
	log.Printf("Received %s /hash", request.Method)
	switch request.Method {
	case "POST":
		if err := request.ParseForm(); err != nil {
			log.Printf("Error - Unable to parse form data %v", err)
			http.Error(response, "Invalid form data", http.StatusBadRequest)
			return
		}

		resultChannel := make(chan persistenceResult)

		go func() {
			dataStore := persistence.GetInstance()
			delay := 5 * time.Second
			id := dataStore.Insert("Placeholder")
			resultChannel <- persistenceResult{
				ID:            id,
				TimeAvailable: time.Now().Add(delay).Format(time.RFC3339),
				URL:           fmt.Sprintf("/hash/%v", id),
			}


			time.Sleep(delay)

			encoded := password.Encode(request.FormValue("password"))
			dataStore.Update(id, encoded)
			log.Println("Persistence of encoded password complete")

		}()

		result := <-resultChannel


		jsonContent, err := json.Marshal(result)

		if err != nil {
			http.Error(response, err.Error(), http.StatusInternalServerError)
			return
		}

		response.Header().Set("Content-Type", "application/json")
		response.WriteHeader(http.StatusAccepted)
		response.Write(jsonContent)
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