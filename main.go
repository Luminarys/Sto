package main

import (
	"bytes"
	"crypto/md5"
	"math/rand"
	"fmt"
	"github.com/hoisie/web"
	"io"
	"os"
    "os/signal"
    "syscall"
	"time"
	"strings"
	"path/filepath"
)

func TempFileName(prefix, suffix string) string {
    s := rand.NewSource(time.Now().UTC().UnixNano())
    r := rand.New(s)
    alphabet := "abcdefghijklmnopqrstuvwxyz0123456789"
    name := ""
    for i := 0; i < 5; i++ {
        idx := r.Intn(len(alphabet))
        name += string(alphabet[idx])
    }
	return filepath.Join(name + suffix)
}

func Md5(r io.Reader) string {
	hash := md5.New()
	io.Copy(hash, r)
	return fmt.Sprintf("%x", hash.Sum(nil))
}

func hello(ctx *web.Context, val string) string {
	ctx.SetHeader("X-Powered-By", "web.go", true)
	ctx.SetHeader("X-Frame-Options", "DENY", true)
	ctx.SetHeader("X-Frame-Options", "DENY", true)
	ctx.SetHeader("Location", "/grill.png", true)
	ctx.SetHeader("Connection", "close", true)
	return "hello " + val
}

func handlePost(ctx *web.Context, val string, updateURL chan<-string) string {
	ctx.Request.ParseMultipartForm(10 * 1024 * 1024)
	form := ctx.Request.MultipartForm
	var output bytes.Buffer

	fileHeader := form.File["file"][0]
	filename := fileHeader.Filename
	file, err := fileHeader.Open()
	if err != nil {
		return err.Error()
	}
    hash := Md5(file)
    //TODO: Parse file properly and get the ext
    name := TempFileName(filename, ".tmp")
    //output.WriteString("orig name: " + filename + " new name: " + TempFileName(filename, ".tmp", Md5(file)) + " hash: " + Md5(file) + "\n")
    output.WriteString("<p>file: " + name + " " + hash + "</p>")
    //Send the URL for updating
    updateURL<-name + ":" + hash
	return output.String()
}

func handleGet(ctx *web.Context, val string, getURL chan<-string, sendURL<-chan string) string {
    getURL<-val
    res := <-sendURL
    if res == "" {
        return "File not found"
    }else{
        return "File at: " + res
    }
}

func urlHandler(getURL<-chan string, sendURL chan<-string, updateURL<-chan string){
    //Read in the CSV, then wait for updates
    urls := make(map[string]string)
    f, err := os.OpenFile("files.csv", os.O_APPEND|os.O_WRONLY, 0660)
    if err != nil {
        fmt.Println("Warning, file could not be opened.")
    }

    defer f.Close()
    for {
        select {
            case read := <-getURL:
                if val, ok := urls[read]; ok {
                    sendURL<-val
                }else{
                    sendURL<-""
                }
            case update := <-updateURL:
                s := strings.Split(update, ":")
                urls[s[0]] = s[1]
                fmt.Println("Updated URLs")
                //Write changes to file
                t := time.Now().UTC()
                nl := s[0] + "," + s[1] + "," + t.Format("2006-01-02 15:04:05") + "\n"
                if _, err := f.WriteString(nl); err != nil {
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

    go urlHandler(getURL, sendURL, updateURL)

    //Le clever. Get around their interface only allowing specific values
    //to be passed by wrapping in a function and sending stuff from there
	web.Post("(/api/upload)", func(ctx *web.Context, val string) string{
        return handlePost(ctx, val, updateURL)
    })
    web.Get("/(.*)", func(ctx *web.Context, val string) string{
        return handleGet(ctx, val, getURL, sendURL)
    })
	web.Run("0.0.0.0:9999")
}
