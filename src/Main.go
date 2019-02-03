package main

import (
	"./handlers"
	"./util"
	"fmt"
	"log"
	"net/http"
	"os"
)

func main() {
	port := ":80"
	root := "."
	index := false
	cors := false
	silent := false
	if util.Contains("--help", &os.Args) != -1 {
		fmt.Println("usage: -p [port] -f [folder] [...options]")
		fmt.Println("--index   enable auto serve index.html")
		fmt.Println("--cors    enable Cross-Origin headers")
		fmt.Println("--silent  suppress logging")

		return
	}
	if util.Contains("-p", &os.Args) != -1 {
		if arg, ok := util.ParseArgs("-p"); !ok {
			panic("argv: invalid argv")
		} else {
			port = ":" + arg
		}
	}
	if util.Contains("-f", &os.Args) != -1 {
		if arg, ok := util.ParseArgs("-f"); !ok {
			panic("argv: invalid argv")
		} else {
			root = arg
		}
	}
	if util.Contains("--index", &os.Args) != -1 {
		index = true
	}
	if util.Contains("--silent", &os.Args) != -1 {
		silent = true
	}
	if util.Contains("--cors", &os.Args) != -1 {
		cors = true
	}
	getHandler := handlers.NewHandler(root, index, cors, silent)
	http.HandleFunc("/", getHandler.Handle)
	fmt.Println("Starting server on port" + port)
	err := http.ListenAndServe(port, nil)
	if err != nil {
		log.Fatal(err)
	}
}

