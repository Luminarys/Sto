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

type urlUpdateMsg struct {
	Name      string
	Extension string
	Hash      string
	Size      int
	Response  chan *Response
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
func handleUpload(ctx *web.Context, updateURL chan<- *urlUpdateMsg) string {
	//TODO: Implemente limits with settings.ini or something
	err := ctx.Request.ParseMultipartForm(100 * 1024 * 1024)
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
		if size > 10*1024*1024 {
			return throwErr(&resFiles)
		}
		//Seek back to beginning
		file.Seek(0, 0)
		if err != nil {
			return err.Error()
		}
		hash := getHash(file)
		ext := filepath.Ext(filename)
		oname := filename[0 : len(filename)-len(ext)]
		//Send the hash and ext for updating
		updateResp := make(chan *Response)
		msg := &urlUpdateMsg{Name: oname, Extension: ext, Hash: hash, Size: int(size), Response: updateResp}
		updateURL <- msg
		resp := <-updateResp
		//Even though this is redundant, it might eventually be useful
		if resp.status == "Failure" {
			fmt.Println("Error for file: ", oname)
			fmt.Println(resp.message)
			return throwErr(&resFiles)
		} else if resp.status == "Duplicate" {
			jFile := jsonFile{Hash: hash, Name: filename, URL: resp.message, Size: int(size)}
			resFiles[idx] = jFile
			//Skip creation for duplicates
			continue
		} else {
			jFile := jsonFile{Hash: hash, Name: filename, URL: resp.message, Size: int(size)}
			resFiles[idx] = jFile
		}

		//If file doesn't already exist, create it
		//Split up files into 3 char prefix and 3 char + ext true file name
		//This should reduce stress on the OS's filesystem
		dir := resp.message[0:3]
		fname := resp.message[3:]
		path := "files/" + dir + "/" + fname
		//If the directory doesn't exist create it
		if !exists("files/" + dir) {
			os.Mkdir("files/"+dir, 0766)
		}
		f, err := os.Create(path)
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
	jsonResp := &jsonResponse{Success: true, Files: resFiles}
	jresp, err := json.Marshal(jsonResp)
	if err != nil {
		fmt.Println(err.Error())
	}
	return string(jresp)
}
