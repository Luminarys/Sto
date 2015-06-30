package main

import (
	"github.com/hoisie/web"
    "flag"
    "os"
)

func main() {
	updateURL := make(chan *urlUpdateMsg)

	go handleDB(updateURL)

	//Get around their func only allowing specific values
	//to be passed by wrapping in a function and sending stuff from there
	web.Post("/api/upload", func(ctx *web.Context) string {
		return handleUpload(ctx, updateURL)
	})
	web.Get("/([a-z0-9]{6}.?[a-z0-9]*)", getFile)
	web.Run("0.0.0.0:8080")
}
