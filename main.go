package main

import (
	"github.com/hoisie/web"
)

func main() {
	getURL := make(chan string)
	getResp := make(chan *Response)
	updateURL := make(chan string)
	updateResp := make(chan *Response)

	go handleDB(getURL, getResp, updateURL, updateResp)

	//Get around their func only allowing specific values
	//to be passed by wrapping in a function and sending stuff from there
	web.Post("/api/upload", func(ctx *web.Context) string {
		return handleUpload(ctx, updateURL, updateResp)
	})
	web.Get("/([a-z0-9]{6}.?[a-z0-9]*)", func(ctx *web.Context, val string) string {
		return getFile(ctx, val, getURL, getResp)
	})
	web.Run("0.0.0.0:8080")
}
