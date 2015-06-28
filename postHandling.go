package main

import (
	"encoding/json"
	"fmt"
	"github.com/hoisie/web"
	"io"
	"os"
	"path/filepath"
)

type jsonFile struct {
	Hash string `json:"hash"`
	Name string `json:"name"`
	URL  string `json:"url"`
	Size int    `json:"size"`
}

type jsonResponse struct {
	Success bool       `json:"success"`
	Files   []jsonFile `json:"files"`
}

func throwErr(arr *[]jsonFile) string {
	jsonResp := &jsonResponse{Success: false, Files: *arr}
	jresp, err := json.Marshal(jsonResp)
	if err != nil {
		fmt.Println(err.Error())
	}
	return string(jresp)
}

//Handles POST upload requests. the updateURL is used to pass messages
//to the urlHandler indicating that the DB should be updated.
func handleUpload(ctx *web.Context, updateURL chan<- string, updateResp <-chan *Response) string {
	//TODO: Implemente limits with settings.ini or something
	err := ctx.Request.ParseMultipartForm(500 * 1024 * 1024)
	if err != nil {
		return "Error handling form!\n"
	}
	form := ctx.Request.MultipartForm

	//Loop through and append to the response struct
	resFiles := make([]jsonFile, len(form.File["files[]"]))
	for idx, fileHeader := range form.File["files[]"] {
		filename := fileHeader.Filename
		file, err := fileHeader.Open()
		size, err := file.Seek(0, 2)
		if err != nil {
			return throwErr(&resFiles)
		}
		if size > 50*1024*1024 {
			return throwErr(&resFiles)
		}
		//Seek back to beginning
		file.Seek(0, 0)
		if err != nil {
			return err.Error()
		}
		hash := Md5(file)
		ext := filepath.Ext(filename)
		//Send the hash and ext for updating
		updateURL <- ext + ":" + hash
		resp := <-updateResp
		//Even though this is redundant, it might eventually be useful
		if resp.status == "Failure" {
			return throwErr(&resFiles)
		} else {
			jFile := jsonFile{Hash: hash, Name: filename, URL: resp.message, Size: int(size)}
			resFiles[idx] = jFile
		}

		//If file doesn't already exist, create it
		if _, err := os.Stat("files/" + hash); os.IsNotExist(err) {
			f, err := os.Create("files/" + hash)
			if err != nil {
				return throwErr(&resFiles)
			}
			_, err = file.Seek(0, 0)
			if err != nil {
				return throwErr(&resFiles)
			}
			_, err = io.Copy(f, file)
			if err != nil {
				return throwErr(&resFiles)
			}
		}
	}
	jsonResp := &jsonResponse{Success: true, Files: resFiles}
	jresp, err := json.Marshal(jsonResp)
	if err != nil {
		fmt.Println(err.Error())
	}
	return string(jresp)
}
