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
	port := ":8080"
	root := "."

	index := false
	cors := false
	silent := false
	auth := false

	if util.Contains("--help", &os.Args) != -1 || util.Contains("-h", &os.Args) != -1 || util.Contains("-help", &os.Args) != -1 || util.Contains("help", &os.Args) != -1 {
		util.PrintHelp()
		os.Exit(0)
	}
	if util.Contains("-p", &os.Args) != -1 {
		if arg, ok := util.ParseArgs("-p"); !ok {
			util.PrintHelp()
			os.Exit(0)
		} else {
			port = ":" + arg
		}
	}
	if util.Contains("-f", &os.Args) != -1 {
		if arg, ok := util.ParseArgs("-f"); !ok {
			util.PrintHelp()
			os.Exit(0)
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
	if util.Contains("--auth", &os.Args) != -1 {
		auth = true
	}
	handler := handlers.NewHandler(root, index, cors, silent, auth)
	if handler.Auth {
		http.HandleFunc("/auth", handler.HandleAuth)
	}
	http.HandleFunc("/", handler.Handle)
	fmt.Println("Starting server on port" + port)
	err := http.ListenAndServe(port, nil)
	if err != nil {
		log.Fatal(err)
	}
}
