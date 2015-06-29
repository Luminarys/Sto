package main

import (
	"github.com/hoisie/web"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

//Handles file retrieval. It uses getURL to send a hash to the urlHandler, and listens on
//sendURL for the proper filename.
func getFile(ctx *web.Context, file string) string {
    dir := file[0:3]
    fname := file[3:]
    path := "files/" + dir + "/" + fname
	//Open the file
	f, err := os.Open(path)
	if err != nil {
		return "Error reading file!\n"
	}

	//Get MIME
	r, err := ioutil.ReadAll(f)
	if err != nil {
		return "Error reading file!\n"
	}
	mime := http.DetectContentType(r)

	_, err = f.Seek(0, 0)
	if err != nil {
		return "Error reading the file\n"
	}
	//This is weird - ServeContent supposedly handles MIME setting
	//But the Webgo content setter needs to be used too
	//In addition, ServeFile doesn't work, ServeContent has to be used
	ctx.ContentType(mime)
	http.ServeContent(ctx.ResponseWriter, ctx.Request, "files/"+fname, time.Now(), f)
	return ""
}
