OUT=out
NAME=goserve
MAIN=src/Main.go
FLAGS=-a -v -ldflags '-w -extldflags "-static"'
COMPILER_OPTS=CGO_ENABLED=0 GOOS=linux GOARCH=amd64
GO_PATH=${GOPATH}

default: build

install: build
	sudo cp ./$(OUT)/$(NAME) /usr/bin/

dep:
	go mod tidy

run:
	go run $(MAIN)

.PHONY: build
build: $(MAIN) dep
	mkdir -p $(OUT)
	$(COMPILER_OPTS) go build $(FLAGS) -o $(OUT)/$(NAME) $(MAIN)

