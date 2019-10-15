OUT=out
NAME=goserve
MAIN=src/Main.go
FLAGS=-a -v -ldflags '-w -extldflags "-static"'

default: build

install:
	sudo ln -sf $(shell pwd)/$(OUT)/$(NAME) /usr/bin/$(NAME)

dep:
	go get github.com/dgrijalva/jwt-go
	go get github.com/go-ini/ini
	go get github.com/gorilla/websocket

build: $(MAIN)	
	mkdir -p $(OUT)
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build $(FLAGS) -o $(OUT)/$(NAME) $(MAIN)

