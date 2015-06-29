package main

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
)

type Response struct {
	status  string
	message string
}

//Handles DB requests by using channels with select to lock access to operations. This ensures that
//files.csv stays updated and maps URLs to hashes(the actual file names).
//The response struct is used to handle responses
func handleDB(updateURL <-chan *urlUpdateMsg) {
	//Create DB if it doesn't exist
	if !exists("./sqlite.db") {
		success := createDB()
		if !success {
			panic("Fatal Error, DB could not be created")
		}
	}
	//Intialize the DB
	db, err := sql.Open("sqlite3", "./sqlite.db")
	if err != nil {
		panic("Fatal Error, DB could not be opened")
	}
	defer db.Close()

	//List of banned extensions
	bannedExts := [30]string{".ade", ".adp", ".bat", ".chm", ".cmd", ".com", ".cpl", ".exe", ".hta", ".ins", ".isp", ".jse", ".lib", ".lnk", ".mde", ".msc", ".msp", ".mst", ".pif", ".scr", ".sct", ".shb", ".sys", ".vb", ".vbe", ".vbs", ".vxd", ".wsc", ".wsf", ".wsh"}
	for {
		select {
		case updateUrlsReq := <-updateURL:
			//Block this operation -- Needs testing to see if the sqlite library
			//fully supports write concurrency
			updateURLs(db, updateUrlsReq, &bannedExts)
		}
	}
}

func updateURLs(db *sql.DB, req *urlUpdateMsg, bannedExts *[30]string) {
	ext := req.Extension
	hash := req.Hash
	origName := req.Name
	size := req.Size
	respChan := req.Response
	//Check if the hash is already present in the DB
	//Should this be a prepared statement? We generate the hash, so it's unlikely
	//that this is exploitable, but it is a possibility
	rows, err := db.Query("SELECT name FROM files WHERE hash = '" + hash + "'")
	if rows.Next() {
		var res string
		rows.Scan(&res)
		respChan <- &Response{status: "Duplicate", message: res}
		return
	}

	for _, e := range *bannedExts {
		if ext == e {
			respChan <- &Response{status: "Failure", message: "This extension is banned, please try again!"}
			return
		}
	}
	//updateURLs(ext, hash)
	//Generate random names until an available slot is there - This might need
	//to be capped, as it could take a LONG time
	name := ""
	for name = RandFileName(ext); exists("files/" + name[0:3] + "/" + name[3:]); name = RandFileName(ext) {
	}

	tx, err := db.Begin()
	if err != nil {
		respChan <- &Response{status: "Failure", message: "Database not functioning properly!"}
		return
	}

	stmt, err := tx.Prepare("Insert into files(name, hash, origname, size) values(?, ?, ?, ?)")
	if err != nil {
		respChan <- &Response{status: "Failure", message: "Database not functioning properly!"}
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(name, hash, origName, size)
	if err != nil {
		respChan <- &Response{status: "Failure", message: "Database transaction failed!"}
		return
	}

	tx.Commit()
	if err != nil {
		respChan <- &Response{status: "Failure", message: "Database transaction failed!"}
		return
	}
	respChan <- &Response{status: "Success", message: name}
}
