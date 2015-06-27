package main

import (
	"github.com/hoisie/web"
	"io"
	"os"
	"path/filepath"
)

//Handles POST upload requests. the updateURL is used to pass messages
//to the urlHandler indicating that the DB should be updated.
func handleUpload(ctx *web.Context, updateURL chan<- string, updateResp <-chan *Response) string {
	//TODO: Implemente limits with settings.ini or something
	err := ctx.Request.ParseMultipartForm(50 * 1024 * 1024)
	if err != nil {
		return "Error handling form!\n"
	}
	form := ctx.Request.MultipartForm

	fileHeader := form.File["file"][0]
	filename := fileHeader.Filename
	file, err := fileHeader.Open()
	size, err := file.Seek(0, 2)
	if err != nil {
		return "Error parsing file!\n"
	}
	if size > 50*1024*1024 {
		return "File too big!\n"
	}
	//Seek back to beginning
	file.Seek(0, 0)
	if err != nil {
		return err.Error()
	}
	hash := Md5(file)
	//If file doesn't already exist, create it
	if _, err := os.Stat("files/" + hash); os.IsNotExist(err) {
		f, err := os.Create("files/" + hash)
		if err != nil {
			return "Error, file could not be created.\n"
		}
		_, err = file.Seek(0, 0)
		if err != nil {
			return "Error reading the file\n"
		}
		_, err = io.Copy(f, file)
		if err != nil {
			return "Error, file could not be written to.\n"
		}
	}

	ext := filepath.Ext(filename)
	//Send the hash and ext for updating
	updateURL <- ext + ":" + hash
	resp := <-updateResp
	//Even though this is redundant, it might eventually be useful
	if resp.status == "Failure" {
		return resp.message + "\n"
	} else {
		return resp.message + "\n"
	}
}