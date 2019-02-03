package util

import (
	"fmt"
	"os"
	"path/filepath"
)

func GenerateHTML(flist *[]os.FileInfo, rel string) []byte {
	upper := "<!DOCTYPE html>\n" +
		"<html lang=\"en\">\n" +
		"<head>\n" +
		"<style>" +
		"*{font-family:Ubuntu,Calibri}" +
		"a{text-decoration:none;font-size:20px;}" +
		"a:hover{color:red!important;}" +
		"td{width:300px;overflow:hidden;}" +
		"</style>" +
		"<meta charset=\"UTF-8\">\n" +
		"<meta name=\"viewport\" content=\"width=device-width, initial-scale=1.0\">\n" +
		"<meta http-equiv=\"X-UA-Compatible\" content=\"ie=edge\">\n" +
		"<title></title>\n" +
		"</head>\n"
	lower :=
		"<body>\n" +
		"</body>\n" +
		"</html>"
	back := fmt.Sprintf("<a style=\"text-decoration:underline;color:black;\" href=\"%s\">&larr;%s</a><br><table><tr><th>Filename</th><th>Size</th></tr>", filepath.Dir(filepath.Dir(rel)), "Back")
	upper += back
	for _, f := range sortDir(flist) {
		upper += "<tr>"
		if f.IsDir() {
			upper += fmt.Sprintf("<td><a style=\"color:black;\" href=\"%s/\">&#128193;%s/</a></td><td style=\"text-align:center\">DIR</td>", rel + f.Name(),f.Name())
		} else {
			upper += fmt.Sprintf("<td><a style=\"color:blue;\" href=\"%s\">&#128462;%s</a></td><td style=\"text-align:center\">%s</td>", rel + f.Name(), f.Name(), byteCountBinary(f.Size()))
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

func ParseArgs(q string) (string, bool){
	index := Contains(q, &os.Args)
	if index == -1 {
		return "", false
	}
	if len(os.Args) == index + 1 {
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

func ContainsFile(q string, dir *[]os.FileInfo) bool{
	for _, file := range *dir {
		if file.Name() == q {
			return true
		}
	}
	return false
}
