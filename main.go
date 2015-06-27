package main

import (
	"bytes"
	"crypto/md5"
	"encoding/csv"
	"fmt"
	"github.com/hoisie/web"
	"io"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

//Generates a random file name and appends on the provided
//extension.
func RandFileName(ext string) string {
	s := rand.NewSource(time.Now().UTC().UnixNano())
	r := rand.New(s)
	alphabet := "abcdefghijklmnopqrstuvwxyz0123456789"
	name := ""
	for i := 0; i < 5; i++ {
		idx := r.Intn(len(alphabet))
		name += string(alphabet[idx])
	}
	return filepath.Join(name + ext)
}

//Returns MD5 hash of a provided file
func Md5(r io.Reader) string {
	hash := md5.New()
	io.Copy(hash, r)
	return fmt.Sprintf("%x", hash.Sum(nil))
}

//Handles POST upload requests. the updateURL is used to pass messages
//to the urlHandler indicating that the DB should be updated.
func handleUpload(ctx *web.Context, updateURL chan<- string) string {
	//TODO: Implemente limits with settings.ini or something
	err := ctx.Request.ParseMultipartForm(50 * 1024 * 1024)
	if err != nil {
		return "Error handling form!\n"
	}
	form := ctx.Request.MultipartForm
	var output bytes.Buffer

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

	extension := filepath.Ext(filename)
	name := RandFileName(extension)
	output.WriteString("name: " + name)
	//Send the URL for updating
	updateURL <- name + ":" + hash
	return output.String() + "\n"
}

//Handles file retrieval. It uses getURL to send a hash to the urlHandler, and listens on
//sendURL for the proper filename.
func getFile(ctx *web.Context, val string, getURL chan<- string, sendURL <-chan string) string {
	getURL <- val
	res := <-sendURL
	if res == "" {
		return "File not found\n"
	} else {
		//Open the file
		f, err := os.Open("files/" + res)
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
		http.ServeContent(ctx.ResponseWriter, ctx.Request, "files/"+res, time.Now(), f)
		return ""
	}
}

//Handles URL processing by using channels with select to lock access to operations. This ensures that
//files.csv stays updated and maps URLs to hashes(the actual file names)
func handleURLs(getURL <-chan string, sendURL chan<- string, updateURL <-chan string) {
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

	for {
		select {
		case read := <-getURL:
			if val, ok := urls[read]; ok {
				sendURL <- val
			} else {
				sendURL <- ""
			}
		case update := <-updateURL:
			//TODO - Verify that there isn't an existing entry in the map
			s := strings.Split(update, ":")
			urls[s[0]] = s[1]
			fmt.Println("Updated URLs")
			//Write changes to file with timestamp for convenience
			t := time.Now().UTC()
			nl := s[0] + "," + s[1] + "," + t.Format("2006-01-02 15:04:05") + "\n"
			if _, err := fout.WriteString(nl); err != nil {
				fmt.Println("Warning, file could not be written to.")
				fmt.Println(err.Error())
			}
		}
	}
}

func main() {
	getURL := make(chan string)
	sendURL := make(chan string)
	updateURL := make(chan string)

	go handleURLs(getURL, sendURL, updateURL)

	//Get around their func only allowing specific values
	//to be passed by wrapping in a function and sending stuff from there
	web.Post("/api/upload", func(ctx *web.Context) string {
		return handleUpload(ctx, updateURL)
	})
	web.Get("/(.*)", func(ctx *web.Context, val string) string {
		return getFile(ctx, val, getURL, sendURL)
	})
	web.Run("0.0.0.0:8080")
}
