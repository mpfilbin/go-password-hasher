package server

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

type RequestHandler func(http.ResponseWriter, *http.Request)

type ApplicationServer struct {
	serveMux *http.ServeMux
	stats *Statistics
}

func NewAppServer() *ApplicationServer{
	return &ApplicationServer{
		serveMux: http.NewServeMux(),
		stats: &Statistics{},
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

func (server *ApplicationServer) Listen(port int) {

	server.RegisterHandler("/shutdown", server.shutdown)
	server.RegisterHandler("/stats", server.reportStatistics)

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