package main

import (
	"github.com/mpfilbin/go-password-hasher/src/http/server"
	"github.com/mpfilbin/go-password-hasher/src/password"
	"github.com/mpfilbin/go-password-hasher/src/statistics"
)

func main() {
	application := server.NewAppServer()

	application.RegisterHandler("/stats", statistics.TrackTiming(statistics.Report))
	application.RegisterHandler("/hash", statistics.TrackTiming(password.EncodeAndPersist))
	application.RegisterHandler("/hash/", statistics.TrackTiming(password.LookupEncodingByID))

	application.Listen(8001)
}
