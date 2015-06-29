package main

import (
	"crypto/sha1"
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"io"
	"math/rand"
	"os"
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

//Returns hash of a provided file
func getHash(r io.Reader) string {
	hash := sha1.New()
	io.Copy(hash, r)
	return fmt.Sprintf("%x", hash.Sum(nil))
}

func exists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return true
}

func createDB() bool {
	db, err := sql.Open("sqlite3", "./sqlite.db")
	if err != nil {
		return false
	}
	defer db.Close()

	sqlStmt := `
    CREATE TABLE files 
    (id integer primary key, name text, hash text, origname text, size integer, Timestamp DATETIME DEFAULT CURRENT_TIMESTAMP);
    `
	_, err = db.Exec(sqlStmt)
	if err != nil {
		return false
	}
	return true
}
