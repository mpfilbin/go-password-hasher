package server

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/mpfilbin/go-password-hasher/src/password"
	"github.com/mpfilbin/go-password-hasher/src/persistence"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"time"
)

type RequestHandler func(http.ResponseWriter, *http.Request)

type persistenceResult struct {
	TimeAvailable string `json:"timeAvailable"`
	URL           string `json:"url"`

}

type ApplicationServer struct {
	serveMux *http.ServeMux
	stats *Statistics
	dataStore *persistence.Repository
}

func NewAppServer() *ApplicationServer{
	return &ApplicationServer{
		serveMux: http.NewServeMux(),
		stats: &Statistics{},
		dataStore: persistence.NewRepository(),
	}
}

func (server *ApplicationServer) RegisterHandler(route string, handler RequestHandler) {
	server.serveMux.HandleFunc(route, func(response http.ResponseWriter, request *http.Request){
		log.Printf("Received %v %v", request.Method, request.URL.Path)
		server.stats.IncrementRequestCount()

		start := time.Now().UnixNano()
		handler(response, request)
		stop := time.Now().UnixNano()

		duration := (stop - start)/int64(time.Microsecond)
		server.stats.AddDuration(duration)
		log.Printf("Request handled in %d microseconds", duration)
	})
}

func (server *ApplicationServer) shutdown(response http.ResponseWriter, request *http.Request) {
	if request.Method == "GET" {
		process, err := os.FindProcess(os.Getpid())
		if err != nil {
			log.Printf("Unable to terminate process due to error: %v", err.Error())
			http.Error(response, err.Error(), http.StatusInternalServerError)
		}
		process.Signal(os.Interrupt)
		return
	}
	http.NotFound(response, request)
}

func (server *ApplicationServer) reportStatistics(response http.ResponseWriter, request *http.Request) {
	if request.Method == http.MethodGet {
		server.stats.UpdateAverageRequestDuration()
		jsonContent, err := json.Marshal(server.stats)

		if err != nil {
			http.Error(response, err.Error(), http.StatusInternalServerError)
			return
		}

		response.Header().Set("Content-Type", "application/json")
		response.Write(jsonContent)
		return
	}
	http.NotFound(response, request)
}

func (server *ApplicationServer) LookupEncodingByID(response http.ResponseWriter, request *http.Request) {
	if request.Method == http.MethodGet {
		idSegment := strings.TrimPrefix(request.URL.Path, "/hash/")
		id, err := strconv.ParseInt(idSegment, 10, 64)

		if err != nil {
			log.Printf("Error - Unable to parse ID from URL\n %v", request.URL.Path)
			http.Error(response, err.Error(), http.StatusBadRequest)
		}

		encodedHash, err := server.dataStore.Get(uint64(id))
		response.Write([]byte(encodedHash))
		return
	}
	http.NotFound(response, request)
}

func (server *ApplicationServer) EncodeAndPersist(response http.ResponseWriter, request *http.Request) {
	if request.Method == http.MethodPost {
		if err := request.ParseForm(); err != nil {
			log.Printf("Error - Unable to parse form data %v", err)
			http.Error(response, "Invalid form data", http.StatusBadRequest)
			return
		}

		resultChannel := make(chan persistenceResult)
		go func() {
			delay := 5 * time.Second
			id := server.dataStore.Insert("")
			resultChannel <- persistenceResult{
				TimeAvailable: time.Now().Add(delay).Format(time.RFC3339),
				URL:           fmt.Sprintf("/hash/%v", id),
			}


			time.Sleep(delay)

			encoded := password.Encode(request.FormValue("password"))
			server.dataStore.Update(id, encoded)
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
		return
	}
	http.NotFound(response, request)
}

func (server *ApplicationServer) Listen(port int) {
	server.RegisterHandler("/shutdown", server.shutdown)
	server.RegisterHandler("/stats", server.reportStatistics)
	server.RegisterHandler("/hash", server.EncodeAndPersist)
	server.RegisterHandler("/hash/", server.LookupEncodingByID)

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)

	address := fmt.Sprintf(":%d", port)
	httpServer := &http.Server{
		Addr:address,
		Handler: server.serveMux,
	}

	log.Printf("Listening on %v", address)
	go func() {
		if err := httpServer.ListenAndServe(); err != nil {
			log.Fatal(err.Error())
		}
	}()

	<-stop
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	httpServer.Shutdown(ctx)
}