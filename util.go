package main

import (
	"crypto/sha1"
	"fmt"
	"io"
	"math/rand"
	"path/filepath"
	"time"
    "os"
)

//Generates a random file name and appends on the provided
//extension.
func RandFileName(ext string) string {
	s := rand.NewSource(time.Now().UTC().UnixNano())
	r := rand.New(s)
	alphabet := "abcdefghijklmnopqrstuvwxyz0123456789"
	name := ""
	for i := 0; i < 6; i++ {
		idx := r.Intn(len(alphabet))
		name += string(alphabet[idx])
	}
	return filepath.Join(name + ext)
}

//Returns hash of a provided file
func getHash(r io.Reader) string {
	hash := sha1.New()
	io.Copy(hash, r)
	return fmt.Sprintf("%x", hash.Sum(nil))
}

func exists(path string) (bool) {
    _, err := os.Stat(path)
    if err == nil { return true}
    if os.IsNotExist(err) { return false}
    return true
}
