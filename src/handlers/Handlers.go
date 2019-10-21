package handlers

import (
	"../livereload"
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
	Root       string
	Index      bool
	Cors       bool
	Silent     bool
	Auth       bool
	LiveReload bool
	password   string
	secret     []byte
}

type HandlerOptions struct {
	Root       string
	Index      bool
	Cors       bool
	Silent     bool
	Auth       bool
	LiveReload bool
	Password   string
}

func NewHandler(opt *HandlerOptions) *Handler {
	// update path
	wd := ""
	if !path.IsAbs(opt.Root) {
		wd, _ = os.Getwd()
		opt.Root = strings.Replace(opt.Root, "/", string(filepath.Separator), -1)
		wd = path.Join(wd, opt.Root)
	} else {
		wd = opt.Root
	}
	wd = strings.Replace(wd, "/", string(filepath.Separator), -1)
	fmt.Println("path: " + wd)
	if _, err := os.Stat(wd); err != nil {
		log.Fatal("path: invalid source path")
	}
	// end update path
	// parse auth data
	handler := &Handler{wd,
		opt.Index,
		opt.Cors,
		opt.Silent,
		opt.Auth,
		opt.LiveReload,
		opt.Password,
		[]byte("d3c9e23120f8849f9e7f8132fbe5400757440493ae11789bfeacc5eabba33e95")}
	if opt.Auth {
		handler.updateAuthParams()
	}
	// end parse auth data
	return handler
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
					return h.secret, nil
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
					if h.LiveReload {
						page, err := livereload.InjectLiveReload(absP + "/" + "index.html")
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
						w.Header().Set("Content-Type", "text/html; charset=utf-8")
						http.ServeFile(w, r, path.Join(absP, "index.html"))
					}
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
				if h.LiveReload {
					page, err := livereload.InjectLiveReload(absP + "/" + fi.Name())
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
					} else {
						http.ServeFile(w, r, absP)
					}
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
				if psh == h.password {
					expires := time.Now().Unix() + int64(24*time.Hour)
					token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{ExpiresAt: expires, Issuer: "server"})
					tokenString, _ := token.SignedString([]byte(h.secret))
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

func getHash(str string) string {
	h := sha256.New()
	h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil))
}

func (h *Handler) updateAuthParams() {
	if i, err := ini.Load("server.ini"); err == nil {
		pw := i.Section("auth").Key("password").String()
		h.password = getHash(pw)
		if util.Contains("-s", &os.Args) == -1 {
			h.secret = []byte(i.Section("auth").Key("secret").String())
		}
	} else {
		fmt.Printf("auth: default password is '%s'\n", h.password)
		h.password = getHash(h.password)
	}
}
