package handlers

import (
	"../util"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/go-ini/ini"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

type Handler struct {
	Root   string
	Index  bool
	Cors   bool
	Silent bool
	Auth   bool
}

//default values
var password = "admin"
var secret = []byte("d3c9e23120f8849f9e7f8132fbe5400757440493ae11789bfeacc5eabba33e95")

func NewHandler(root string, index bool, cors bool, silent bool, auth bool) *Handler {
	// update path
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
	// end update path
	// parse auth data
	if auth {
		updateAuthParams()
	}
	// end parse auth data
	return &Handler{wd, index, cors, silent, auth}
}
func (h *Handler) Handle(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		if h.Cors {
			w.Header().Set("Access-Control-Allow-Origin", "*")
		}
		if !h.Silent {
			fmt.Println(r.Host+"\t", r.URL.String())
		}
		if h.Auth {
			if token, err := r.Cookie("Authorization"); err != nil {
				fmt.Println(err)
				http.Redirect(w, r, "/auth", 301)
				return
			} else {
				tokenString := strings.Split(token.Value, "Bearer ")[1]
				if _, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
					if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
						return nil, fmt.Errorf("jwt: unexpected signing method %v", token.Header["alg"])
					}
					return secret, nil
				}); err != nil {
					fmt.Println(err)
					http.Redirect(w, r, "/auth", 307)
					return
				}

			}
		}
		// AUTH GOOD
		p := r.URL.String()
		p = strings.Replace(p, "/", string(filepath.Separator), -1)
		absP := ""
		if p == "/" || p == "\\" {
			absP = h.Root
		} else {
			absP = path.Join(h.Root, p)
		}
		if fi, err := os.Stat(absP); err == nil && fi.IsDir() {
			if dir, err := ioutil.ReadDir(absP); err == nil {
				if util.ContainsFile("index.html", &dir) && h.Index {
					page, err := util.InjectLiveReload(absP + "/" + "index.html")
					if err != nil {
						fmt.Println(err)
						w.WriteHeader(http.StatusInternalServerError)
						if _, err := fmt.Fprint(w, "( ͠° ͟ʖ ͡°) 500 INTERNAL SERVER ERROR"); err != nil {
							fmt.Println(err)
						}
					}
					w.Header().Set("Content-Type", "text/html; charset=utf-8")
					w.Header().Set("Content-Length", strconv.Itoa(len(page)))
					if _, err := w.Write(page); err != nil {
						fmt.Println(err)
					}
					//w.Header().Set("Content-Type", "text/html; charset=utf-8")
					//http.ServeFile(w, r, path.Join(absP, "index.html"))
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
				if _, err := fmt.Fprint(w, "( ͠° ͟ʖ ͡°) 500 INTERNAL SERVER ERROR"); err != nil {
					fmt.Println(err)
				}
			}
		} else if err == nil {
			if strings.HasSuffix(strings.ToLower(fi.Name()), ".html") {
				page, err := util.InjectLiveReload(absP + "/" + fi.Name())
				if err != nil {
					fmt.Println(err)
					w.WriteHeader(http.StatusInternalServerError)
					if _, err := fmt.Fprint(w, "( ͠° ͟ʖ ͡°) 500 INTERNAL SERVER ERROR"); err != nil {
						fmt.Println(err)
					}
				}
				w.Header().Set("Content-Type", "text/html; charset=utf-8")
				w.Header().Set("Content-Length", strconv.Itoa(len(page)))
				if _, err := w.Write(page); err != nil {
					fmt.Println(err)
				}
			} else {
				http.ServeFile(w, r, absP)
			}
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
func (h *Handler) HandleAuth(w http.ResponseWriter, r *http.Request) {
	if h.Auth {
		if r.Method == "GET" {
			w.Header().Set("Content-Type", "text/html")
			w.Header().Set("Content-Length", strconv.Itoa(len(util.RenderLoginPage())))
			if _, err := w.Write(util.RenderLoginPage()); err != nil {
				fmt.Println(err)
			}
		} else {
			if err := r.ParseForm(); err != nil {
				w.WriteHeader(http.StatusUnauthorized)
				if _, err := fmt.Fprint(w, "( ͠° ͟ʖ ͡°) 401 UNAUTHORIZED"); err != nil {
					fmt.Println(err)
				}
			} else {
				ps := r.Form.Get("password")
				psh := getHash(ps)
				if psh == password {
					expires := time.Now().Unix() + int64(24*time.Hour)
					token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{ExpiresAt: expires, Issuer: "server"})
					tokenString, _ := token.SignedString([]byte(secret))
					cookie := http.Cookie{Name: "Authorization", Value: fmt.Sprintf("Bearer %s", tokenString), Path: "/", Expires: time.Now().Add(24 * time.Hour)}
					http.SetCookie(w, &cookie)
					http.Redirect(w, r, "/", 301)
				} else {
					http.Redirect(w, r, "/auth", 301)
				}
			}
		}
	} else {
		w.WriteHeader(http.StatusNotImplemented)
		if _, err := fmt.Fprint(w, "( ͠° ͟ʖ ͡°) 501 NOT IMPLEMENTED"); err != nil {
			fmt.Println(err)
		}
	}
}
func SetPassword(p string) {
	hm := sha256.New()
	hm.Write([]byte(p))
	password = hex.EncodeToString(hm.Sum(nil))
}
func SetSecret(s string) {
	secret = []byte(s)
}

func getHash(str string) string {
	h := sha256.New()
	h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil))
}

func updateAuthParams() {
	if util.Contains("-pw", &os.Args) != -1 {
		if arg, ok := util.ParseArgs("-pw"); !ok {
			util.PrintHelp()
			os.Exit(0)
		} else {
			password = getHash(arg)
		}
	} else if i, err := ini.Load("server.ini"); err == nil {
		pw := i.Section("auth").Key("password").String()
		password = getHash(pw)
		if util.Contains("-s", &os.Args) == -1 {
			secret = []byte(i.Section("auth").Key("secret").String())
		}
	} else {
		password = getHash(password)
		fmt.Println("auth: no -pw option or server.ini file")
		fmt.Println("auth: default password is 'admin'")
	}
	if util.Contains("-s", &os.Args) != -1 {
		if arg, ok := util.ParseArgs("-s"); !ok {
			util.PrintHelp()
			os.Exit(0)
		} else {
			secret = []byte(arg)
		}
	}
}
