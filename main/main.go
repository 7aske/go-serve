package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

func main() {
	port := ":80"
	if len(os.Args) == 2 {
		port = ":" + os.Args[1]
	}
	http.HandleFunc("/", handler)
	createDirIfNotExist("uploads")
	fmt.Println("Starting server on port" + port)
	err := http.ListenAndServe(port, nil)
	if err != nil {
		log.Fatal(err)
	}
}

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.Method)
	if err := r.ParseForm(); err == nil {
		if r.Method == "GET" {
			if r.URL.String() == "/" {
				http.ServeFile(w, r, "static/index.html")
			} else {
				http.ServeFile(w, r, "static"+r.URL.String())
			}
		} else if r.Method == "POST" {
			filename := r.Form["filename"][0]
			text := r.Form["text"][0]
			if text != "" && filename != "" {
				saveFile(filename, text)
				if _, err := fmt.Fprint(w, "File saved successfully"); err != nil {
					fmt.Println(err)
				}
			} else {
				http.Redirect(w, r, "/", http.StatusSeeOther)
			}
		} else {
			w.WriteHeader(404)
			if _, err := fmt.Fprint(w, "404 NOT FOUND"); err != nil {
				fmt.Println(err)
			}
		}
	} else {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		if _,err :=	fmt.Fprint(w, "Internal Server Error"); err != nil {
			fmt.Println(err)
		}
	}
}

func saveFile(filename string, text string) {
	f, err := os.Create("uploads/" + filename)
	if err != nil {
		fmt.Println(err)
	}

	l, err := f.WriteString(text)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(l, "bytes written")
	err = f.Close()
	if err != nil {
		fmt.Println(err)
	}
}
func createDirIfNotExist(dir string) {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err = os.MkdirAll(dir, 0755)
		if err != nil {
			panic(err)
		}
	}
}
