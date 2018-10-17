package main

import (
	"github.com/mpfilbin/go-password-hasher/http/server"
)

func main() {
	application := server.NewAppServer()
	application.Listen(8080)
}
