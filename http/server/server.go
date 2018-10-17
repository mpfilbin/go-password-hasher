package server

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/mpfilbin/go-password-hasher/password"
	"github.com/mpfilbin/go-password-hasher/persistence"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"time"
)

const (
	ContentTypeJson      = "application/json; charset=utf-8"
	ContentTypePlaintext = "text/plain; charset=utf-8"
)

type requestHandler func(http.ResponseWriter, *http.Request)

type encodingResult struct {
	TimeAvailable string `json:"timeAvailable"`
	URL           string `json:"url"`
}

// ApplicationServer provides an HTTP interface to enable the encoding and retrieval of passwords
type ApplicationServer struct {
	serveMux  *http.ServeMux
	stats     *statistics
	dataStore *persistence.Repository
}

// NewAppServer instantiates a new instance of ApplicationServer
func NewAppServer() *ApplicationServer {
	return &ApplicationServer{
		serveMux:  http.NewServeMux(),
		stats:     &statistics{},
		dataStore: persistence.NewRepository(),
	}
}

func (server *ApplicationServer) registerHandler(route string, handler requestHandler) {
	server.serveMux.HandleFunc(route, func(response http.ResponseWriter, request *http.Request) {
		log.Printf("Received %v %v", request.Method, request.URL.Path)
		server.stats.IncrementRequestCount()

		start := time.Now().UnixNano()
		handler(response, request)
		stop := time.Now().UnixNano()

		duration := (stop - start) / int64(time.Microsecond)
		server.stats.AddDuration(duration)
		log.Printf("Request handled in %d microseconds", duration)
	})
}

func (server *ApplicationServer) shutdown(response http.ResponseWriter, request *http.Request) {
	if request.Method == http.MethodGet {
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

		response.Header().Set("Content-Type", ContentTypeJson)
		response.Write(jsonContent)
		return
	}
	http.NotFound(response, request)
}

func (server *ApplicationServer) lookupEncodingByID(response http.ResponseWriter, request *http.Request) {
	if request.Method == http.MethodGet {
		idSegment := strings.TrimPrefix(request.URL.Path, "/hash/")
		id, err := strconv.ParseInt(idSegment, 10, 64)

		if err != nil {
			log.Printf("Error - Unable to parse ID from URL\n %v", request.URL.Path)
			http.Error(response, err.Error(), http.StatusBadRequest)
		}

		encodedHash, err := server.dataStore.Get(int64(id))
		response.Write([]byte(encodedHash))
		return
	}
	http.NotFound(response, request)
}

func (server *ApplicationServer) encodeAndPersist(response http.ResponseWriter, request *http.Request) {
	if request.Method == http.MethodPost {
		if err := request.ParseForm(); err != nil {
			log.Printf("Error - Unable to parse form data %v", err)
			http.Error(response, "Invalid form data", http.StatusBadRequest)
			return
		}

		resultChannel := make(chan encodingResult)
		go func() {
			delay := 5 * time.Second
			id := server.dataStore.Insert("")
			resultChannel <- encodingResult{
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

		response.Header().Set("Content-Type", ContentTypeJson)
		response.WriteHeader(http.StatusAccepted)
		response.Write(jsonContent)
		return
	}
	http.NotFound(response, request)
}

// Listen binds the instance of ApplicationServer to a network port so that it may receive HTTP traffic
func (server *ApplicationServer) Listen(port int) {
	server.registerHandler("/shutdown", server.shutdown)
	server.registerHandler("/stats", server.reportStatistics)
	server.registerHandler("/hash", server.encodeAndPersist)
	server.registerHandler("/hash/", server.lookupEncodingByID)

	idleConnsClosed := make(chan struct{})

	httpServer := &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: server.serveMux,
	}

	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt)
		<-sigint
		if err := httpServer.Shutdown(context.Background()); err != nil {
			// Error from closing listeners, or context timeout:
			log.Printf("HTTP server Shutdown: %v", err)
		}
	}()

	log.Printf("Listening on %v", httpServer.Addr)
	if err := httpServer.ListenAndServe(); err != nil {
		log.Fatal(err.Error())
	}

	<-idleConnsClosed
}
