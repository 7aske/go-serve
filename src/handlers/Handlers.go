package handlers

import (
	"../util"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
)

type GetHandler struct {
	Root   string
	Index  bool
	Cors   bool
	Silent bool
}

func NewHandler(root string, index bool, cors bool, silent bool) *GetHandler {
	wd := ""
	if !path.IsAbs(root) {
		wd, _ = os.Getwd()
		root = strings.Replace(root, "/", string(filepath.Separator), -1)
		wd = path.Join(wd, root)
	} else {
		wd = root
	}
	wd = strings.Replace(wd, "/", string(filepath.Separator), -1)
	fmt.Println("path: " + wd)
	if _, err := os.Stat(wd); err != nil {
		log.Fatal("path: invalid source path")
	}
	return &GetHandler{wd, index, cors, silent}
}
func (h *GetHandler) Handle(w http.ResponseWriter, r *http.Request) {
	if !h.Silent {
		fmt.Println(r.Method, r.URL.String(), r.Host)
	}
	if r.Method == "GET" {
		if h.Cors {
			w.Header().Set("Access-Control-Allow-Origin", "*")
		}
		p := r.URL.String()
		p = strings.Replace(p, "/", string(filepath.Separator), -1)
		absP := ""
		if p == "/" || p == "\\" {
			fmt.Println(p)
			absP = h.Root
		} else {
			absP = path.Join(h.Root, p)
		}
		fmt.Println(absP)
		if fi, err := os.Stat(absP); err == nil && fi.IsDir() {
			if dir, err := ioutil.ReadDir(absP); err == nil {
				if util.ContainsFile("index.html", &dir) && h.Index {
					w.Header().Set("Content-Type", "text/html; charset=utf-8")
					http.ServeFile(w, r, path.Join(absP, "index.html"))
				} else {
					page := util.GenerateHTML(&dir, r.URL.String())
					w.Header().Set("Content-Type", "text/html; charset=utf-8")
					w.Header().Set("Content-Length", strconv.Itoa(len(page)))
					if _, err := w.Write(page); err != nil {
						fmt.Println(err)
					}
				}
			} else {
				fmt.Println(err)
				w.WriteHeader(http.StatusInternalServerError)
				if _, err := fmt.Fprint(w, "Internal Server Error"); err != nil {
					fmt.Println(err)
				}
			}
		} else if err == nil {
			http.ServeFile(w, r, absP)
		} else {
			w.WriteHeader(http.StatusNotFound)
			if _, err := fmt.Fprint(w, "( ͠° ͟ʖ ͡°) 404 NOT FOUND"); err != nil {
				fmt.Println(err)
			}
		}
	} else {
		w.WriteHeader(http.StatusNotImplemented)
		if _, err := fmt.Fprint(w, "( ͠° ͟ʖ ͡°) 501 NOT IMPLEMENTED"); err != nil {
			fmt.Println(err)
		}
	}

}
