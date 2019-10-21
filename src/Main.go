package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"./handlers"
	"./livereload"
)

func PrintHelp() {
	_, _ = fmt.Fprintln(os.Stderr, "usage: goserve [dir] [...flags]")
	_, _ = fmt.Fprintln(os.Stderr, "")
	_, _ = fmt.Fprintln(os.Stderr, "--port, -p <port>       specify server port")
	_, _ = fmt.Fprintln(os.Stderr, "--dir, -d <folder>      specify source folder")
	_, _ = fmt.Fprintln(os.Stderr, "--index, -i             enable auto serve index.html")
	_, _ = fmt.Fprintln(os.Stderr, "--reload, -r            enable auto-reloading of html pages on fs changes")
	_, _ = fmt.Fprintln(os.Stderr, "--cors                  enable Cross-Origin headers")
	_, _ = fmt.Fprintln(os.Stderr, "--silent, -s            suppress logging")
	_, _ = fmt.Fprintln(os.Stderr, "--auth, -a              enables authentication")
	_, _ = fmt.Fprintln(os.Stderr, "--password, -pw <pass>  specify auth password")
	_, _ = fmt.Fprintln(os.Stderr, "example:")
	_, _ = fmt.Fprintln(os.Stderr, "")
	_, _ = fmt.Fprintln(os.Stderr, "goserve ./static --index --reload")
}

func main() {
	port := flag.String("port", "8080", "specify server port")
	flag.StringVar(port, "p", *port, "specify server port")

	dir := flag.String("dir", ".", "specify source folder")
	flag.StringVar(dir, "d", *dir, "specify source folder")

	auth := flag.Bool("auth", false, "enables authentication")
	flag.BoolVar(auth, "a", *auth, "enables authentication")

	cors := flag.Bool("cors", false, "enable Cross-Origin headers")

	silent := flag.Bool("silent", false, "suppress logging")
	flag.BoolVar(silent, "s", *silent, "suppress logging")

	index := flag.Bool("index", false, "enable auto serve index.html")
	flag.BoolVar(index, "i", *index, "enable auto serve index.html")

	reload := flag.Bool("reload", false, "enable auto-reloading of html pages on fs changes")
	flag.BoolVar(reload, "r", *reload, "enable auto-reloading of html pages on fs changes")

	password := flag.String("password", "admin", "specify auth password")
	flag.StringVar(password, "pw", *password, "specify auth password")

	flag.Usage = func() {
		PrintHelp()
	}
	flag.Parse()
	if len(flag.Args()) > 0 {
		if flag.Arg(0) != *dir {
			*dir = flag.Arg(0)
		}
	}

	if !strings.HasPrefix(*port, ":") {
		*port = ":" + *port
	}

	if *index && *reload {
		go livereload.ListenAndServe()
	} else {
		*reload = false
	}
	handlerOptions := handlers.HandlerOptions{
		Root:       *dir,
		Index:      *index,
		Cors:       *cors,
		Silent:     *silent,
		Auth:       *auth,
		LiveReload: *reload,
		Password:   *password,
	}
	handler := handlers.NewHandler(&handlerOptions)
	if handler.Auth {
		http.HandleFunc("/auth", func(w http.ResponseWriter, r *http.Request) {
			handler.HandleAuth(w, r)
		})
	}
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		handler.Handle(w, r)
	})
	fmt.Println("Starting server on port " + strings.TrimPrefix(*port, ":"))
	log.Fatal(http.ListenAndServe(*port, nil))
}
