package main

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
)

//Generic response template that will be sent to functions requesting DB operations
type Response struct {
	status  string
	message string
}

type writeMsg struct {
	name     string
	origName string
	hash     string
	size     int
}

//Handles DB requests by using channels with select to lock access to operations. This ensures that
//files.csv stays updated and maps URLs to hashes(the actual file names).
//The response struct is used to handle responses
func handleDB(updateURL <-chan *urlUpdateMsg) {
	//Intialize the DB
	db, err := sql.Open("sqlite3", "./sqlite.db")
	if err != nil {
		panic("Fatal Error, DB could not be opened")
	}
	defer db.Close()

	writeData := make(chan *writeMsg, 10)

	//List of banned extensions
	bannedExts := [30]string{".ade", ".adp", ".bat", ".chm", ".cmd", ".com", ".cpl", ".exe", ".hta", ".ins", ".isp", ".jse", ".lib", ".lnk", ".mde", ".msc", ".msp", ".mst", ".pif", ".scr", ".sct", ".shb", ".sys", ".vb", ".vbe", ".vbs", ".vxd", ".wsc", ".wsf", ".wsh"}
	for {
		select {
		case updateUrlsReq := <-updateURL:
			//Since this operation will only do DB reads, we can make it concurrent
			go updateURLs(db, updateUrlsReq, &bannedExts, writeData)
		case writeReq := <-writeData:
			//Block this operation since it involves actual writes
			writeToDB(db, writeReq)
		}
	}
}

func updateURLs(db *sql.DB, req *urlUpdateMsg, bannedExts *[30]string, writeReq chan<- *writeMsg) {
	ext := req.Extension
	hash := req.Hash
	origName := req.Name
	size := req.Size
	respChan := req.Response

	//Make sure that the extension is valid
	for _, e := range *bannedExts {
		if ext == e {
			respChan <- &Response{status: "Failure", message: "This extension is banned, please try again!"}
			return
		}
	}
	//Check if the hash is already present in the DB
	//Should this be a prepared statement? We generate the hash, so it's unlikely
	//that this is exploitable, but it is a possibility
	row, err := db.Query("SELECT name FROM files WHERE hash = '" + hash + "'")
	defer row.Close()

	if err != nil {
		respChan <- &Response{status: "Failure", message: "Could not query DB!"}
		return
	}

	if row.Next() {
		var res string
		row.Scan(&res)
		respChan <- &Response{status: "Duplicate", message: res}
		return
	}

	//updateURLs(ext, hash)
	//Generate random names until an available slot is there - This might need
	//to be capped, as it could take a LONG time
	name := ""
	for name = RandFileName(ext); exists("files/" + name[0:3] + "/" + name[3:]); name = RandFileName(ext) {
	}
	//The channel is buffered so this should be responded to almost instantly
	writeReq <- &writeMsg{name: name, origName: origName, hash: hash, size: size}
	respChan <- &Response{status: "Success", message: name}
}

func writeToDB(db *sql.DB, info *writeMsg) {
	tx, err := db.Begin()
	if err != nil {
		//respChan <- &Response{status: "Failure", message: "Could not initiate transaction!"}
		return
	}

	stmt, err := tx.Prepare("Insert into files(name, hash, origname, size) values(?, ?, ?, ?)")
	if err != nil {
		//respChan <- &Response{status: "Failure", message: "Could not create prepared statement!"}
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(info.name, info.hash, info.origName, info.size)
	if err != nil {
		//respChan <- &Response{status: "Failure", message: "Could not execute prepared statement!"}
		return
	}

	tx.Commit()
	if err != nil {
		//respChan <- &Response{status: "Failure", message: "Could not commit transaction!"}
		return
	}
}
