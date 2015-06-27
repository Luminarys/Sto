package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"strings"
	"time"
)

type Response struct {
	status  string
	message string
}

//Handles DB requests by using channels with select to lock access to operations. This ensures that
//files.csv stays updated and maps URLs to hashes(the actual file names).
//The response struct is used to handle responses
func handleDB(getURL <-chan string, getResp chan<- *Response, updateURL <-chan string, updateResp chan<- *Response) {
	//Read in the CSV, then wait for updates
	urls := make(map[string]string)
	if _, err := os.Stat("files.csv"); os.IsNotExist(err) {
		f, err := os.Create("files.csv")
		if err != nil {
			panic("Fatal Error, files.csv could not be created")
		}
		f.Close()
	}
	fin, err := os.Open("files.csv")
	if err != nil {
		panic("Fatal Error, files.csv could not be opened.")
	}

	reader := csv.NewReader(fin)
	data, err := reader.ReadAll()
	if err != nil {
		panic("Fatal Error, files.csv is not formatted properly")
	}

	for _, col := range data {
		urls[col[0]] = col[1]
	}

	fin.Close()

	fout, err := os.OpenFile("files.csv", os.O_APPEND|os.O_WRONLY, 0660)
	if err != nil {
		panic("Fatal Error, files.csv could not be opened.")
	}

	defer fout.Close()

	//List of banned extensions
	bannedExts := [30]string{".ade", ".adp", ".bat", ".chm", ".cmd", ".com", ".cpl", ".exe", ".hta", ".ins", ".isp", ".jse", ".lib", ".lnk", ".mde", ".msc", ".msp", ".mst", ".pif", ".scr", ".sct", ".shb", ".sys", ".vb", ".vbe", ".vbs", ".vxd", ".wsc", ".wsf", ".wsh"}

	for {
		select {
		case read := <-getURL:
			if val, ok := urls[read]; ok {
				getResp <- &Response{status: "Success", message: val}
			} else {
				getResp <- &Response{status: "Failure", message: "File not found!"}
			}
		case update := <-updateURL:
			s := strings.Split(update, ":")
			ext := s[0]
			for _, e := range bannedExts {
				if ext == e {
					updateResp <- &Response{status: "Failure", message: "This extension is banned, please try again!"}
				}
			}
			hash := s[1]

			//Verify that the generated name is unique
			name := RandFileName(ext)
			//notExists should be true
			for _, notExists := urls[name]; notExists; name = RandFileName(ext) {
				_, notExists = urls[name]
			}
			urls[name] = hash
			fmt.Println("Updated URLs")
			//Write changes to file with timestamp for convenience
			t := time.Now().UTC()
			nl := s[0] + "," + s[1] + "," + t.Format("2006-01-02 15:04:05") + "\n"
			if _, err := fout.WriteString(nl); err != nil {
				updateResp <- &Response{status: "Failure", message: "Warning, file could not be recorded!"}
			} else {
				updateResp <- &Response{status: "Success", message: name}
			}
		}
	}
}
