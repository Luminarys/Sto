build: check-deps
	go build -o Sto src/main.go src/dbHandler.go src/getHandling.go src/postHandling.go src/util.go

check-deps:
	if test -z "$$GOPATH"; then \
	GOPATH=~/Gopath; \
	mkdir -p ~/Gopath; \
	echo "GOPATH undefined, changed to $$GOPATH"; \
	GOPATH=$$GOPATH go get github.com/hoisie/web; \
	GOPATH=$$GOPATH go get github.com/mattn/go-sqlite3; \
	else \
	go get github.com/hoisie/web; \
	go get github.com/mattn/go-sqlite3; \
	fi
