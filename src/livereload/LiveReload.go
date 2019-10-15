package livereload

import (
	"bufio"
	"fmt"
	"github.com/gorilla/websocket"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

var conn *websocket.Conn

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func wsEndpoint(w http.ResponseWriter, r *http.Request) {
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
	}
	fmt.Println("livereload connection established")
	conn = ws
}


func ListenAndServe() {
	http.HandleFunc("/ws", wsEndpoint)
	ln, err := net.Listen("tcp", ":33900")
	if err != nil {
		log.Fatal(err)
	}
	go watch()
	log.Fatal(http.Serve(ln, nil))
}

func watch() {
	var nodes []os.FileInfo
	var sum int64
	cwd, err := os.Getwd()
	root, err := os.Open(cwd)
	if err != nil {
		fmt.Print(err.Error())
	}
	err = filepath.Walk(root.Name(), func(path string, info os.FileInfo, err error) error {
		if !strings.HasPrefix(info.Name(),"."){
			nodes = append(nodes, info)
		}
		return nil
	})
	if err != nil {
		fmt.Print(err.Error())
	}
	for _, node := range nodes {
		sum += node.ModTime().Unix()
	}
	for {
		var newnodes []os.FileInfo
		var newsum int64
		err = filepath.Walk(root.Name(), func(path string, info os.FileInfo, err error) error {
			if !strings.HasPrefix(info.Name(),"."){
				newnodes = append(newnodes, info)
			}
			return nil
		})
		if err != nil {
			fmt.Println(err)
		}
		for _, node := range newnodes {
			newsum += node.ModTime().Unix()
		}
		if newsum != sum {
			sum = newsum
			if conn != nil {
				messageType, _, err := conn.ReadMessage()
				if err != nil {
					log.Println(err)
				}
				if err := conn.WriteMessage(messageType, []byte("reload")); err != nil {
					log.Println(err)
				}
			}
		}
		time.Sleep(time.Millisecond * 1000)
	}
}

func InjectLiveReload(pth string) ([]byte, error) {
	file, err := os.Open(pth)
	if err != nil {
		return []byte{}, err
	}
	reader := bufio.NewReader(file)
	fileContents, err := ioutil.ReadAll(reader)
	if err != nil {
		return []byte{}, err
	}
	return appendScript(fileContents), nil
}

func appendScript(buffer []byte) []byte {
	var script = "<script>var protocol = window.location.protocol === 'http:' ? 'ws://' : 'wss://';var address = protocol + window.location.hostname + ':33900' + '/ws';var socket = new WebSocket(address);socket.onopen = function() {socket.send('connected');};socket.onmessage = function (msg) { if (msg.data === 'reload') location.reload();};</script>"
	return append(buffer, []byte(script)...)
}
