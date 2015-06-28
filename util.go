package main

import (
	"crypto/md5"
	"fmt"
	"io"
	"math/rand"
	"path/filepath"
	"time"
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

//Returns MD5 hash of a provided file
func Md5(r io.Reader) string {
	hash := md5.New()
	io.Copy(hash, r)
	return fmt.Sprintf("%x", hash.Sum(nil))
}
