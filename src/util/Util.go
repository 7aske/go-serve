package util

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

func RenderLoginPage() []byte {
	return []byte(`<!doctype html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport"
          content="width=device-width, user-scalable=no, initial-scale=1.0, maximum-scale=1.0, minimum-scale=1.0">
    <meta http-equiv="X-UA-Compatible" content="ie=edge">
    <title>Deployment Server</title>
	<style>
		*{font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, "Helvetica Neue", Arial, sans-serif, "Apple Color Emoji", "Segoe UI Emoji", "Segoe UI Symbol", "Noto Color Emoji";}
		body,html {width: 100%;overflow-x: hidden;background-color: #222222;color: #ffffff;}
		input {margin-left: -8px;border-radius: 8px;padding: 10px;border: 2px solid #666666;font-size: 24px;}
	</style>
</head>
<body style="text-align: center;">
	<h1>Admin Login</h1>
	<form method="POST" action="/auth">
		<input type="password" name="password" placeholder="Password" />
	</form>
</body>
</html>`)
}

func GenerateHTML(flist *[]os.FileInfo, rel string) []byte {
	upper := `<!DOCTYPE html><html lang="en">
		<head>
		<style>
		*{font-family:Ubuntu,Calibri}
		a{text-decoration:none;font-size:20px;}
		a:hover{color:red!important;}
		tr:hover td {color:red!important;}
		td{width:500px;overflow:hidden;}
		td+td{width:200px;}
		</style>
		<meta charset="UTF-8">
		<meta name="viewport" content="width=device-width, initial-scale=1.0">
		<meta http-equiv="X-UA-Compatible" content="ie=edge">
		<title></title>
		</head>`
	upper += fmt.Sprintf(`<a style="font-size:24px;color:black;" href="/">%s</a><hr>`, rel)
	lower := `<body></body></html>`
	back := fmt.Sprintf(`<a style="text-decoration:underline;color:black;" href="%s">&larr;%s</a><br><br><table><tr><th style="text-align:left;">Filename</th><th>Size</th></tr>`, filepath.Dir(filepath.Dir(rel)), "Back")
	upper += back
	for _, f := range sortDir(flist) {
		upper += "<tr>"
		if f.IsDir() {
			upper += fmt.Sprintf(`<td><a style="color:black;" href="%s/">&#128193;%s/</a></td><td style="text-align:center">DIR</td>`, rel+f.Name(), f.Name())
		} else {
			if f.Name() != "server.ini" {
				upper += fmt.Sprintf(`<td><a style="color:blue;" href="%s">&#128462;%s</a></td><td style="text-align:center">%s</td>`, rel+f.Name(), f.Name(), byteCountBinary(f.Size()))
			}
		}
		upper += "</tr>"
	}
	upper += lower
	return []byte(upper)

}

func byteCountBinary(b int64) string {
	const unit = 1024
	if b < unit {
		return fmt.Sprintf("%d B", b)
	}
	div, exp := int64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %ciB", float64(b)/float64(div), "KMGTPE"[exp])
}

func sortDir(list *[]os.FileInfo) []os.FileInfo {
	var folders []os.FileInfo
	var files []os.FileInfo
	for _, f := range *list {
		if f.IsDir() {
			folders = append(folders, f)
		} else {
			files = append(files, f)
		}
	}
	return append(folders, files...)
}

func ParseArgs(q string) (string, bool) {
	index := Contains(q, &os.Args)
	if index == -1 {
		return "", false
	}
	if len(os.Args) == index+1 {
		return "", false
	}
	return os.Args[index+1], true
}
func Contains(q string, s *[]string) int {
	for i, str := range *s {
		if str == q {
			return i
		}
	}
	return -1
}

func ContainsFile(q string, dir *[]os.FileInfo) bool {
	for _, file := range *dir {
		if file.Name() == q {
			return true
		}
	}
	return false
}
func PrintHelp() {
	fmt.Println("usage: [...options] [...flags]")
	fmt.Println()
	fmt.Println("-p  <port>      specify server port")
	fmt.Println("-f  <folder>    specify source folder")
	fmt.Println("-pw <pass>      specify server password")
	fmt.Println("                only works with --auth flag")
	fmt.Println("-s  <secret>    specify hash secret")
	fmt.Println("                only works with --auth flag")
	fmt.Println("--index         enable auto serve index.html")
	fmt.Println("--cors          enable Cross-Origin headers")
	fmt.Println("--silent        suppress logging")
	fmt.Println("--auth          enables authentication")
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
	var script = "<script>var protocol = window.location.protocol === 'http:' ? 'ws://' : 'wss://';var address = protocol + window.location.hostname + ':33900' + window.location.pathname + '/ws';var socket = new WebSocket(address);socket.onmessage = function (msg) {if (msg.data == 'reload') window.location.reload();};</script>"
	return append(buffer, []byte(script)...)
}
