package main

import (
	"fmt"
	"net/http"
	"os"
)

func main() {
	http.HandleFunc("/", handler)
	createDirIfNotExist("uploads")
	fmt.Println("Server started on port 80.")
	http.ListenAndServe(":80", nil)
}

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.Method)
	r.ParseForm()
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
 			fmt.Fprint(w, "File saved successfully")
		} else {
			http.Redirect(w, r, "/", http.StatusSeeOther)
		}
	} else {
		w.WriteHeader(404)
		fmt.Fprint(w, "404 NOT FOUND")
	}

}

func saveFile(filename string, text string) {
	f, err := os.Create("uploads/" + filename)
	if err != nil {
		return
	}

	l, err := f.WriteString(text)
	if err != nil {
		return
	}

	fmt.Println(l, "bytes written")
	err = f.Close()
	if err != nil {
		return
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