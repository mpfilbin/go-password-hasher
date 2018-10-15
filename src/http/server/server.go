package server

import (
	"context"
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
}


func NewAppServer() *ApplicationServer{
	return &ApplicationServer{serveMux: http.NewServeMux()}
}


func (server *ApplicationServer) RegisterHandler(route string, handler RequestHandler) {
	server.serveMux.HandleFunc(route, handler)
}

func (server *ApplicationServer) shutdown(response http.ResponseWriter, request *http.Request) {
	if request.Method == "GET" {
		process, err := os.FindProcess(os.Getpid())
		if err != nil {
			log.Printf("Unable to terminate process due to error: %v", err.Error())
			http.Error(response, err.Error(), http.StatusInternalServerError)
		}
		process.Signal(os.Interrupt)
	} else {
		http.NotFound(response, request)
	}
}

func (server *ApplicationServer) Listen(port int) {

	server.RegisterHandler("/shutdown", server.shutdown)

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