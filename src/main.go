package main

import (
	"flag"
	"github.com/hoisie/web"
	"github.com/mattn/go-session-manager"
	"log"
	"os"
	"runtime"
)

var manager = session.NewSessionManager(logger)
var logger = log.New(os.Stdout, "", log.Ldate|log.Ltime)

type User struct {
	Username string
	Password string
}

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

	manager.OnStart(func(session *session.Session) {
		logger.Printf("Start session(\"%s\")", session.Id)
	})
	manager.OnEnd(func(session *session.Session) {
		logger.Printf("End session(\"%s\")", session.Id)
	})

	web.Config.CookieSecret = "ayyyLmao"

	updateURL := make(chan *urlUpdateMsg)
	login := make(chan *loginReq)
	go handleDB(updateURL, login)

	//Set route handlers
	web.Post("/api/upload", func(ctx *web.Context) string {
		return handleUpload(ctx, updateURL)
	})
	web.Post("/login", func(ctx *web.Context) string {
		return handleLogin(ctx, login)
	})
	web.Get("/([a-z0-9]{6}.?[a-z0-9]*)", getFile)
}

func main() {
	port := flag.String("port", "8080", "The port that Sto should listen on. By default it is 8080.")
	procs := flag.Int("procs", 1, "The maximum number of processes that can be used by Go. The default value is one, but at least two are recommended in order to maximize performance.")
	flag.Parse()

	if *procs > runtime.NumCPU() {
		panic("Fatal error: You tried to use more processes than there are CPUs available")
	}
	runtime.GOMAXPROCS(*procs)
	startUp()

	web.Run("0.0.0.0:" + *port)
}
