package livereload

import (
	"fmt"
	"net"
	"os"
	"path/filepath"
	"time"
)

var conn net.Conn

func Test() {
	ln, err := net.Listen("tcp", ":33900")
	if err != nil {
		fmt.Println(err)
	}
	go watch()
	for {
		conn, err = ln.Accept()
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println("received conn")
	}
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
		nodes = append(nodes, info)
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
			newnodes = append(newnodes, info)
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
			fmt.Println(sum)
			fmt.Println(newsum)
			if conn != nil {
				_, _ = conn.Write([]byte("reload\n"))
				_ = conn.Close()
			}
		}
		time.Sleep(time.Millisecond * 1000)
	}
}
