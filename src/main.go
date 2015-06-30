package main

import (
	"github.com/hoisie/web"
    "flag"
    "os"
)

func startUp() {
	//Create DB if it doesn't exist
	if !exists("./sqlite.db") {
		success := createDB()
		if !success {
			panic("Fatal Error, DB could not be created")
		}
	}
    //Create Files/ if it doesn't exist
    if !exists("./files/") {
        os.Mkdir("files/", 0766)
    }

    //Set route handlers
	updateURL := make(chan *urlUpdateMsg)
	go handleDB(updateURL)

	web.Post("/api/upload", func(ctx *web.Context) string {
		return handleUpload(ctx, updateURL)
	})
	web.Get("/([a-z0-9]{6}.?[a-z0-9]*)", getFile)
}

func main() {
    port := flag.String("p", "8080", "The port that Sto should listen on. By default it is 8080")
    flag.Parse()

    startUp()

    web.Run("0.0.0.0:" + *port)
}
