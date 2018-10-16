package main

import (
	"github.com/mpfilbin/go-password-hasher/src/http/server"
)

func main() {
	application := server.NewAppServer()
	application.Listen(8001)
}
