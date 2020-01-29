package zip

import (
	"archive/zip"
	"fmt"
	"io/ioutil"
	"os"
	"path"
)

func Folder(baseFolder string) (string, error) {
	tempDir := path.Join(os.TempDir(), "go-serve")
	err := os.MkdirAll(tempDir, 0775)
	if err != nil {
		return "", err
	}

	// If the folder has already been zipped we just return its filepath
	outFilePath := path.Join(tempDir, path.Base(baseFolder)+".zip")
	if _, err := os.Stat(outFilePath); err == nil {
		return outFilePath, nil
	}
	outFile, err := os.Create(outFilePath)
	if err != nil {
		return "", err
	}

	w := zip.NewWriter(outFile)

	addFiles(w, baseFolder, "")

	err = w.Close()
	if err != nil {
		return "", err
	}
	return outFilePath, err
}

func addFiles(w *zip.Writer, basePath, baseInZip string) {
	// Open the Directory
	files, err := ioutil.ReadDir(basePath)
	if err != nil {
		fmt.Println(err)
	}

	for _, file := range files {
		if !file.IsDir() {
			dat, err := ioutil.ReadFile(path.Join(basePath, file.Name()))
			if err != nil {
				fmt.Println(err)
			}

			// Add some files to the archive.
			f, err := w.Create(path.Join(baseInZip, file.Name()))
			if err != nil {
				fmt.Println(err)
			}
			_, err = f.Write(dat)
			if err != nil {
				fmt.Println(err)
			}
		} else if file.IsDir() {

			newBase := path.Join(basePath, file.Name(), "/")
			addFiles(w, newBase, path.Join(baseInZip, file.Name(), "/"))
		}
	}
}
