package main

import (
	"flag"
	"github.com/hoisie/web"
	"os"
    "runtime"
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

	updateURL := make(chan *urlUpdateMsg)
	go handleDB(updateURL)

	//Set route handlers
	web.Post("/api/upload", func(ctx *web.Context) string {
		return handleUpload(ctx, updateURL)
	})
	web.Get("/([a-z0-9]{6}.?[a-z0-9]*)", getFile)
}

func main() {
	port := flag.String("port", "8080", "The port that Sto should listen on. By default it is 8080")
	procs := flag.Int("procs", 1, "The maximum number of processes that can be used by Go")
	flag.Parse()

    if *procs > runtime.NumCPU() {
        panic("Fatal error: You tried to use more processes than there are CPUs available")
    }
    runtime.GOMAXPROCS(*procs)
	startUp()

	web.Run("0.0.0.0:" + *port)
}
