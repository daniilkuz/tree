package main

import (
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
)

func printDir(out io.Writer, pathName string, file fs.DirEntry, prefix string) error {
	fileName := file.Name()
	fileInfo, err := os.Stat(filepath.Join(pathName, fileName))
	if err != nil {
		return err
	}

	var stringToWrite string

	if !file.IsDir() {
		var postfix string
		if fileInfo.Size() == 0 {
			postfix = "(empty)"
		} else {
			postfix = fmt.Sprintf("(%db)", fileInfo.Size())
		}
		stringToWrite = fmt.Sprintf("%s%s %s\n", prefix, fileName, postfix)
	} else {
		stringToWrite = prefix + fileName + "\n"
	}
	out.Write([]byte(stringToWrite))

	return nil
}

func buildTree(out io.Writer, prefix string, path string, printFiles bool) error {

	files, _ := os.ReadDir(path)

	if !printFiles {
		var newFiles = []fs.DirEntry{}
		for _, file := range files {
			if file.IsDir() {
				newFiles = append(newFiles, file)
			}
		}
		files = newFiles
	}

	for i, file := range files {
		var newPrefix string
		if i == len(files)-1 {
			printDir(out, path, file, prefix+"└───")
			newPrefix = prefix + "	"
		} else {
			printDir(out, path, file, prefix+"├───")
			newPrefix = prefix + "│	"
		}
		if file.IsDir() {
			buildTree(out, newPrefix, filepath.Join(path, file.Name()), printFiles)
		}
	}
	return nil
}

func dirTree(out io.Writer, path string, printFiles bool) error {
	buildTree(out, "", path, printFiles)
	return nil
}

func main() {
	out := os.Stdout
	if !(len(os.Args) == 2 || len(os.Args) == 3) {
		panic("usage go run main.go . [-f]")
	}
	path := os.Args[1]
	printFiles := len(os.Args) == 3 && os.Args[2] == "-f"
	err := dirTree(out, path, printFiles)
	if err != nil {
		panic(err.Error())
	}
}
