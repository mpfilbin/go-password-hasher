package main

import (
	"github.com/mpfilbin/go-password-hasher/src/http/server"
	"github.com/mpfilbin/go-password-hasher/src/password"
)

func main() {
	application := server.NewAppServer()

	application.RegisterHandler("/hash", password.EncodeAndPersist)
	application.RegisterHandler("/hash/", password.LookupEncodingByID)

	application.Listen(8001)
}
