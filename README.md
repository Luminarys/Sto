# Sto
A simple file uploader and sharer, based off of Pomf.se. Sto's backend is written in Go, and is designed to be both performant and portable.

## Requirements
Sto requires Go, and the packages `github.com/hoisie/web` and `github.com/mattn/go-sqlite3`, as well as a valid SQLite3 installation which was compiled with enabled concurrency options.

## Setup and Running
Either run `go run *.go` or build the go files and run the resulting executable. You will also want to create a Files directory in the cloned or downloaded github directory prior to execution. A sqlite DB will be automatically created.
To configure Nginx, the following location block should be utilized: 
```
location / {
    proxy_pass http://127.0.0.1:8080;
    proxy_set_header Host $host;
    proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
}
```

## Performance
A benchmark using JMeter and attempting to upload 10 files at a time yielded the following results:
### Sto
![alt text](https://fuwa.se/2z0hl.png/Go_Results.png "Sto Benchmark 1")
### Pomf
![alt text](https://fuwa.se/o3lpw.png/Pomf_Results.png "Pomf Benchmark 1")

While these results should not be treated as completely conclusive, they do provide some indication of the strength of Go as a potential for servers. It should also be noted that Sto utilizes SQLite(to make portability easy) which is also slower than MySQL, which is used by Pomf.

## TODO
* Improve DB Handling
* General efficiency improvements?
* Improve Frontend
* Makefile for optional compilation and minification of assets
