package main

import (
	"crypto/sha1"
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"hash"
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

	// 2 channels: used to give green light for reading into buffer b1 or b2
	readch1, readch2 := make(chan int, 1), make(chan int, 1)

	// 2 channels: used to give green light for hashing the content of b1 or b2
	hashch1, hashch2 := make(chan int, 1), make(chan int, 1)

	// Start signal: Allow b1 to be read and hashed
	readch1 <- 1
	hashch1 <- 1

	go hashHelper(r, hash, readch1, readch2, hashch1, hashch2)

	hashHelper(r, hash, readch2, readch1, hashch2, hashch1)

	return fmt.Sprintf("%x", hash.Sum(nil))
}

func hashHelper(r io.Reader, h hash.Hash, mayRead <-chan int, readDone chan<- int, mayHash <-chan int, hashDone chan<- int) {
	for b, hasMore := make([]byte, 64<<10), true; hasMore; {
		<-mayRead
		n, err := r.Read(b)
		if err != nil {
			if err == io.EOF {
				hasMore = false
			} else {
				panic(err)
			}
		}
		readDone <- 1

		<-mayHash
		_, err = h.Write(b[:n])
		if err != nil {
			panic(err)
		}
		hashDone <- 1
	}
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
