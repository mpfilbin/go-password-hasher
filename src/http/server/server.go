package server

import (
	"fmt"
	"github.com/mpfilbin/go-password-hasher/src/password"
	"github.com/mpfilbin/go-password-hasher/src/statistics"
	"log"
	"net/http"
)


func Listen(port int) {
	serverMux := http.NewServeMux() // Create isolated server mutex

	serverMux.HandleFunc("/hash", statistics.TrackTiming(password.EncodeAndPersist))
	serverMux.HandleFunc("/hash/", statistics.TrackTiming(password.LookupEncodingByID))
	serverMux.HandleFunc("/stats", statistics.TrackTiming(statistics.Report))


	address :=fmt.Sprintf(":%d", port)
	log.Printf("Listening at %s", address)

	err := http.ListenAndServe(address, serverMux) // set listen port
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}