package password

import (
	"encoding/json"
	"fmt"
	"github.com/mpfilbin/go-password-hasher/src/persistence"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type persistenceResult struct {
	ID            uint64 `json:"id"`
	TimeAvailable string `json:"timeAvailable"`
	URL           string `json:"url"`

}


func LookupEncodingByID(response http.ResponseWriter, request *http.Request) {
	switch request.Method {
	case "GET":
		idSegment := strings.TrimPrefix(request.URL.Path, "/hash/")
		id, err := strconv.ParseInt(idSegment, 10, 64)

		if err != nil {
			log.Printf("Error - Unable to parse ID from URL\n %v", request.URL.Path)
			http.Error(response, err.Error(), http.StatusBadRequest)
		}



		dataStore := persistence.GetInstance()
		encodedHash, err := dataStore.Get(uint64(id))
		response.Write([]byte(encodedHash))
	default:
		http.NotFound(response, request)
	}
}

func EncodeAndPersist(response http.ResponseWriter, request *http.Request) {
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
			id := dataStore.Insert("")
			resultChannel <- persistenceResult{
				ID:            id,
				TimeAvailable: time.Now().Add(delay).Format(time.RFC3339),
				URL:           fmt.Sprintf("/hash/%v", id),
			}


			time.Sleep(delay)

			encoded := Encode(request.FormValue("password"))
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